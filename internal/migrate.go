package internal

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDB(connStr string) error {
	m, err := migrate.New("file:sql", connStr)
	if err != nil {
		return err
	}
	err = m.Up()
	switch err {
	case migrate.ErrNoChange:
		return nil
	default:
		return err
	}
}
