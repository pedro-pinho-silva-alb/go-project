package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // Driver PostgreSQL
)

var dbpool *pgxpool.Pool

func ensureDatabaseExists() error {
	// Conectar diretamente ao banco de dados "postgres" (administrativo)
	connStr := "postgres://username:password@postgres-service.default.svc.cluster.local:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error connecting to the admin database: %w", err)
	}
	defer db.Close()

	// Verificar se o banco de dados existe
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'go_database')`
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking database existence: %w", err)
	}

	// Criar o banco de dados se não existir
	if !exists {
		_, err = db.Exec("CREATE DATABASE go_database")
		if err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}
		log.Println("Database 'go_database' created successfully")
	} else {
		log.Println("Database 'go_database' already exists")
	}

	return nil
}

func connectPostgres() (*pgxpool.Pool, error) {
	// Conectar ao banco de dados específico
	connStr := "postgres://user:password@localhost:5432/go_database?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failure connecting to PostgreSQL: %w", err)
	}
	return pool, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)

	// Execute database query
	rows, err := dbpool.Query(context.Background(), "SELECT external_id FROM home")
	if err != nil {
		http.Error(w, "Error looking at database", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var external_id string
		if err := rows.Scan(&external_id); err != nil {
			http.Error(w, "Error processing results", http.StatusInternalServerError)
			log.Println("Error scanning db line:", err)
			return
		}
		fmt.Fprintf(w, "External ID: %s\n", external_id)
	}

	fmt.Fprintf(w, "Thank you!\n")
}

func main() {
	log.Println("Waiting for the database to be available...")
	time.Sleep(10 * time.Second) // Ajuste esse valor conforme necessário
	// Garantir que o banco de dados existe
	err := ensureDatabaseExists()
	if err != nil {
		log.Fatalf("Error ensuring database exists: %v", err)
	}

	// Conectar ao banco de dados
	dbpool, err = connectPostgres()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer dbpool.Close()

	// Configurar e iniciar o servidor HTTP
	http.HandleFunc("/", handler)
	fmt.Println("Server is running on http://localhost:8085")
	http.ListenAndServe(":8085", nil)
}
