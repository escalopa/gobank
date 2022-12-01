package util

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
)

func InitDatabase(config Config) (conn *sql.DB, err error) {
	conn, err = sql.Open(config.Driver, config.ConnectionString)
	if err != nil {
		return
	}

	// Connection retry 5 times
	for i := 0; i < 5; i++ {
		err = conn.Ping()
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
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
