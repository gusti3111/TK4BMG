package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_"github.com/lib/pq" // PostgreSQL driver
)

// DB adalah koneksi database global yang akan digunakan oleh repositories
var DB *sql.DB

// ConnectDB initializes the database connection
func ConnectDB() error {

	// Construct connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// NOTE: Pastikan di .env, DB_HOST diatur ke 'db' (nama service Docker Compose)
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify connection
	err = database.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = database
	log.Println("Database connection (BMG) successful!")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed.")
	}
}
