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
	CreatedTime      time.Time
}

var ErrHashNotfound = errors.New("hash not found")

type Storage interface {
	CreateSecret(secret *SecretData) (string, error)
	GetSecret(secretID string) (*SecretData, error)
	Delete(secretID string) error
}

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
	const sqlStmt = "UPDATE secrets SET expireAfterViews = expireAfterViews - 1 WHERE id = $1 RETURNING secret, expireAfterViews, expireAfterTime, createdTime"
	row := s.db.QueryRow(sqlStmt, secretID)
	var secret SecretData
	switch err := row.Scan(&secret.Secret,
		&secret.ExpireAfterViews,
		&secret.ExpireAfterTime,
		&secret.CreatedTime); err {
	case sql.ErrNoRows:
		return nil, ErrHashNotfound
	case nil:
		return &secret, nil
	default:
		return nil, err
	}
}

func (s *PgStorage) Delete(secretID string) error {
	const sqlStmt = "DELETE FROM secrets WHERE id = $1"
	_, err := s.db.Exec(sqlStmt, secretID)
	return err
}

func (s *PgStorage) Close() {
	err := s.db.Close()
	if err != nil {
		log.Printf("failed to close PgConnection : %v", err)
	}
}
