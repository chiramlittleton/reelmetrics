package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB holds the database connection
var DB *sql.DB

func ConnectDB() {
	var err error
	DB, err = sql.Open("pgx", "postgres://user:password@postgres:5432/reelmetrics_db?sslmode=disable")
	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}
	log.Println("✅ Connected to PostgreSQL")
}
