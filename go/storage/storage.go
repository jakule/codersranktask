package storage

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type SecretData struct {
	Secret           string
	ExpireAfterViews int
	ExpireAfterTime  *time.Time
}

type Storage interface {
	CreateSecret(secret *SecretData) (string, error)
	GetSecret(secretID string) (*SecretData, error)
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

func (s *PgStorage) CreateSecret(secret *SecretData) (string, error) {
	const sqlStmt = `INSERT INTO secrets (secret, expireAfterViews, expireAfterTime)
	VALUES ($1, $2, $3)
	RETURNING id`
	var id string
	err := s.db.QueryRow(sqlStmt, secret.Secret, secret.ExpireAfterViews,
		secret.ExpireAfterTime).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *PgStorage) GetSecret(secretID string) (*SecretData, error) {
	const sqlStmt = "SELECT secret, expireAfterViews, expireAfterTime FROM secrets WHERE id = $1"
	row := s.db.QueryRow(sqlStmt, secretID)
	var secret SecretData
	switch err := row.Scan(&secret.Secret,
		&secret.ExpireAfterViews,
		&secret.ExpireAfterTime); err {
	case sql.ErrNoRows:
		return nil, ErrHashNotfound
	case nil:
		return &secret, nil
	default:
		return nil, err
	}
}

func (s *PgStorage) Close() {
	err := s.db.Close()
	if err != nil {
		log.Printf("failed to close PgConnection : %v", err)
	}
}
