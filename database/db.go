package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"

	_"github.com/lib/pq"
)

func DbIn() (db *sql.DB, err error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load .env %v", err)
		return nil, err
	}
	conStr := os.Getenv("CONNECTION_STRING")
	db, err = sql.Open("postgres", conStr)
	if err != nil {
		log.Fatalf("error in connection string %v", err)
		return nil, err

	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("error while connect database %v", err)
		return nil, err

	}
	return db, nil
}
