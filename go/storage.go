package swagger

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateSecret(secret string) (string, error)
	GetSecret(secretID string) (string, error)
}

var ErrHashNotfound = errors.New("hash not found")

type PgStorage struct {
	db *sql.DB
}

func NewPgStorage(connStr string) (Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PgStorage{db: db}, nil
}

func (s *PgStorage) CreateSecret(secret string) (string, error) {
	stmt := `INSERT INTO secrets (secret)
	VALUES ($1)
	RETURNING id`
	var id string
	err := s.db.QueryRow(stmt, secret).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *PgStorage) GetSecret(secretID string) (string, error) {
	row := s.db.QueryRow("SELECT secret FROM secrets WHERE id = $1", secretID)
	var secret string
	switch err := row.Scan(&secret); err {
	case sql.ErrNoRows:
		return "", ErrHashNotfound
	case nil:
		return secret, nil
	default:
		return "", err
	}
}

func (s *PgStorage) Close() {
	err := s.db.Close()
	if err != nil {
		log.Printf("failed to close PgConnection : %v", err)
	}
}
