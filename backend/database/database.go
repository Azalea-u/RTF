package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

// NewDatabase initializes a new database connection and executes the schema.
func NewDatabase(dbpath string, schemaPath string) *Database {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Fatal("\033[31mError:\033[0m" + " Database connection failed - " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("\033[31mError:\033[0m" + " Database connection failed - " + err.Error())
	}

	schema, err := readSchema(schemaPath)
	if err != nil {
		log.Fatal("\033[31mError:\033[0m" + " Failed to read schema - " + err.Error())
	}

	if _, err := db.ExecContext(ctx, schema); err != nil {
		log.Fatal("\033[31mError:\033[0m" + " Schema execution failed - " + err.Error())
	}

	log.Println("\033[32mSuccess:\033[0m" + " Database connection successful")

	return &Database{DB: db}
}

// readSchema reads the SQL schema from a file.
func readSchema(schemaPath string) (string, error) {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Close closes the database connection.
func (db *Database) Close() error {
	if err := db.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}

// GetDB returns the underlying *sql.DB instance.
func (db *Database) GetDB() *sql.DB {
	return db.DB
}
