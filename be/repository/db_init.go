package repository

import (
	"database/sql"
	"log"
	"os"
)

func InitDB() *sql.DB {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=pgdb port=5432 user=postgres password=yourpassword dbname=expenses sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS blocks (
			id TEXT PRIMARY KEY,
			month TEXT UNIQUE,
			locked BOOLEAN
		)`,
		`CREATE TABLE IF NOT EXISTS members (
			id TEXT PRIMARY KEY,
			block_id TEXT,
			name TEXT,
			ratio FLOAT,
			debt FLOAT,
			FOREIGN KEY (block_id) REFERENCES blocks(id)
		)`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id TEXT PRIMARY KEY,
			block_id TEXT,
			payer TEXT,
			amount FLOAT,
			description TEXT,
			created_at TIMESTAMP,
			ratios JSONB,
			FOREIGN KEY (block_id) REFERENCES blocks(id)
		)`,
		`CREATE TABLE IF NOT EXISTS transaction_details (
			transaction_id TEXT,
			member_id TEXT,
			amount FLOAT,
			PRIMARY KEY (transaction_id, member_id)
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS user_logs (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			method TEXT NOT NULL,
			path TEXT NOT NULL,
			ip_address TEXT,
			user_agent TEXT,
			body TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatal(err)
		}
	}

	return db
}
