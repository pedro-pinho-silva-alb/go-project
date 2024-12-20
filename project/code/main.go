package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"project/code/db"
	"project/code/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool

func main() {
	log.Println("Waiting for the database to be available...")
	time.Sleep(3 * time.Second)

	const (
		adminConnStr = "postgres://username:password@postgres-service.default.svc.cluster.local:5432/postgres?sslmode=disable"
		// adminConnStr = "postgres://username:password@localhost:5432/postgres?sslmode=disable"
		databaseName = "go_database"
		userConnStr  = "postgres://user:password@localhost:5432/go_database?sslmode=disable"
	)

	err := db.EnsureDatabaseExists(adminConnStr, databaseName)
	if err != nil {
		log.Fatalf("Error ensuring database exists: %v", err)
	}

	dbpool, err = db.ConnectPostgres(userConnStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer dbpool.Close()

	http.HandleFunc("/", handlers.HomeHandler(dbpool))

	fmt.Println("Server is running on http://localhost:8085")
	http.ListenAndServe(":8085", nil)
}
