package storage

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPgStorage_CreateSecret(t *testing.T) {
	var (
		expectedQuery = `INSERT INTO secrets (secret, expireAfterViews, expireAfterTime)
	VALUES ($1, $2, $3)
	RETURNING id`
		secretText       = "secretMsg"
		expireAfterViews = 5
		expectedID       = "c7cb197a-de61-4190-8735-17ac5a343826"
	)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("cannot create mock %v", err)
	}
	s := &PgStorage{
		db: db,
	}

	mock.ExpectQuery(expectedQuery).
		WithArgs(secretText, expireAfterViews, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	secret := &SecretData{
		Secret:           secretText,
		ExpireAfterViews: expireAfterViews,
		ExpireAfterTime:  nil,
	}

	got, err := s.CreateSecret(secret)
	if err != nil {
		t.Fatalf("failed to create secret %v", err)
	}

	if got != expectedID {
		t.Fatalf("got %s, expected %s", got, expectedID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPgStorage_GetSecret(t *testing.T) {
	var (
		expectedQuery = `UPDATE secrets SET expireAfterViews = expireAfterViews - 1 WHERE id = $1 RETURNING secret, expireAfterViews, expireAfterTime, createdTime`
		secretID      = "c7cb197a-de61-4190-8735-17ac5a343826"
	)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("cannot create mock %v", err)
	}

	var (
		s = &PgStorage{
			db: db,
		}

		secret = &SecretData{
			Secret:           "secretText",
			ExpireAfterViews: 60,
			ExpireAfterTime:  nil,
			CreatedTime:      time.Now(),
		}

		returnedColumns = []string{"secret", "expireAfterViews",
			"expireAfterTime", "createdTime"}
	)

	mock.ExpectQuery(expectedQuery).
		WithArgs(secretID).
		WillReturnRows(sqlmock.NewRows(returnedColumns).
			AddRow(secret.Secret, secret.ExpireAfterViews,
				secret.ExpireAfterTime, secret.CreatedTime))

	got, err := s.GetSecret(secretID)
	if err != nil {
		t.Fatalf("failed to create secret %v", err)
	}

	if !reflect.DeepEqual(got, secret) {
		t.Fatalf("got %v, expected %v", got, secret)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
