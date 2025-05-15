package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using environment variables: %v", err)
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "postgres://postgres:secret@localhost:5432/makerble?sslmode=disable"
		log.Printf("DB_URL not set in environment, using default: %s", dsn)
	}

	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Database connection established")
}

func Close() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Fatalf("Error closing the database connection: %v", err)
		} else {
			fmt.Println("Database connection closed")
		}
	}
}
