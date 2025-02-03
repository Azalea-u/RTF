package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const schemaFilePath = "./db/schema.sql"

type Database struct {
	DB *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("Database connection established and schema initialized successfully")
	return &Database{DB: db}, nil
}

func initSchema(db *sql.DB) error {
	sqlBytes, err := os.ReadFile(schemaFilePath)
	if err != nil {
		return fmt.Errorf("failed to read schema file (%s): %w", schemaFilePath, err)
	}

	if _, err := db.Exec(string(sqlBytes)); err != nil {
		return fmt.Errorf("failed to execute schema SQL: %w", err)
	}

	return nil
}

func (d *Database) Close() error {
	if err := d.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	log.Println("Database connection closed")
	return nil
}
