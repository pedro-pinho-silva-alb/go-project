package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // Driver PostgreSQL
)

func EnsureDatabaseExists(connStrAdmin, databaseName string) error {
	db, err := sql.Open("postgres", connStrAdmin)
	if err != nil {
		return fmt.Errorf("error connecting to the admin database: %w", err)
	}
	defer db.Close()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`
	err = db.QueryRow(query, databaseName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking database existence: %w", err)
	}

	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", databaseName))
		if err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}
		log.Printf("Database '%s' created successfully\n", databaseName)
	} else {
		log.Printf("Database '%s' already exists\n", databaseName)
	}

	return nil
}

func ConnectPostgres(connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failure connecting to PostgreSQL: %w", err)
	}
	return pool, nil
}
