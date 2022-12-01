package util

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
)

func InitDatabase(config Config) (conn *sql.DB, err error) {
	conn, err = sql.Open(config.Driver, config.ConnectionString)
	if err != nil {
		return
	}

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return &sql.DB{}, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migration",
		config.Name, driver)

	if err != nil {
		return &sql.DB{}, err
	}

	m.Up()

	if err != nil {
		return &sql.DB{}, err
	}

	return
}
