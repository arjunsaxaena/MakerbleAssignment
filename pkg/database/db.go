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
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Warning: .env file not found, using environment variables: %v", err)
	} // Didn't knew godotenv required relative path for .env

	dsn := os.Getenv("DB_URL")

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
