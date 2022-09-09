package dbutils

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/davidklassen/testing-workshop/pkg/retry"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func OpenPostgres(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	if err = retry.Retry(5, time.Second, 2, db.Ping); err != nil {
		return nil, fmt.Errorf("failed to ping postgres server")
	}

	return db, nil
}

func MustOpenPostgres(connString string) *sql.DB {
	db, err := OpenPostgres(connString)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	return db
}

func MigratePostgres(connStr, dbName, migrations string) error {
	db, err := OpenPostgres(connStr)
	if err != nil {
		return fmt.Errorf("failed to open postgres: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to get driver instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrations,
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m.Up()
}
