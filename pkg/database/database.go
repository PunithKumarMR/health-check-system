package database

import (
    "database/sql"
    "fmt"
    "time"
    _ "github.com/go-sql-driver/mysql"
)

// Config holds database configuration
type Config struct {
    Host     string
    Port     string
    User     string
    Password string
    Database string
}

// DB wraps the sql.DB connection
type DB struct {
    *sql.DB
}

// Connect establishes a connection to the database
func Connect(cfg Config) (*DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Set connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    // Test the connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
    return db.DB.Close()
}

// Ping tests the database connection
func (db *DB) Ping() error {
    return db.DB.Ping()
}
