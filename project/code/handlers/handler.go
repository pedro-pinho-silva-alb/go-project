package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HomeHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)

		rows, err := dbpool.Query(context.Background(), "SELECT external_id FROM home")
		if err != nil {
			http.Error(w, "Error looking at database", http.StatusInternalServerError)
			log.Println("Error executing query:", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var externalID string
			if err := rows.Scan(&externalID); err != nil {
				http.Error(w, "Error processing results", http.StatusInternalServerError)
				log.Println("Error scanning db line:", err)
				return
			}
			fmt.Fprintf(w, "External ID: %s\n", externalID)
		}

		fmt.Fprintf(w, "Thank you!\n")
	}
}
