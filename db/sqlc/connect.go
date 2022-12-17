package db

import (
	"database/sql"
	"log"

	"github.com/escalopa/gobank/util"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
)

func InitDatabase(config *util.Config) *sql.DB {
	conn, err := sql.Open(config.Get("DATABASE_DRIVER"), config.Get("DATABASE_URL"))
	if err != nil {
		log.Fatalf("cannot open connection to db, err: %s", err)
	}

	err = migrateDB(conn, config.Get("DATABASE_MIGRATION_PATH"))
	if err != nil {
		log.Fatalf("cannot migrate db, err: %s", err)
	}

	return conn
}

func migrateDB(conn *sql.DB, migrationURL string) error {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationURL,
		"postgres", driver)

	if err != nil {
		return err
	}

	m.Up()

	if err != nil {
		return err
	}

	return nil
}
