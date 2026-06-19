package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	// The blank import registers the PostgreSQL driver with the database/sql package.
	// This enables sql.Open("postgres", ...) to work correctly.
	_ "github.com/lib/pq"
)

// Embedding *sql.DB means DB automatically has all sql.DB methods.
type DB struct {
	*sql.DB
}

func Connect() *DB {
	connStr := os.Getenv("DATABASE_URL")

	var database *sql.DB
	var err error

	// Consider making MAX_ATTEMPT declared elsewhere
	MAX_ATTEMPTS := 10
	for i := 0; i < MAX_ATTEMPTS; i++ {
		// sql.Open does NOT immediately connect; it only creates a connection pool.
		database, err = sql.Open("postgres", connStr)
		if err == nil {
			// Ping actually tests the connection to the database.
			err = database.Ping()
		}
		if err == nil {
			fmt.Println("connected to database")
			// Wrap the sql.DB pointer in our custom DB type.
			return &DB{database}
		}
		fmt.Printf("waiting for database... attempt %d/10\n", i+1)
		time.Sleep(2 * time.Second)
	}

	fmt.Println("could not connect to database")
	os.Exit(1)
	return nil
}

// InsertReturningID runs an INSERT and returns the generated id
func (db *DB) InsertReturningID(query string, args ...interface{}) (int, error) {
	var id int
	err := db.QueryRow(query, args...).Scan(&id)
	return id, err
}

// DeleteByID runs a DELETE and returns an error if nothing was deleted
func (db *DB) DeleteByID(query string, id int) error {
	result, err := db.Exec(query, id)
	if err != nil { return err }

	rowsAffected, err := result.RowsAffected()
	if err != nil { return err }
	
	if rowsAffected == 0 {
		return fmt.Errorf("record %d not found", id)
	}
	return nil
}
// GetByID runs a SELECT for a single row and scans into dest
func (db *DB) GetByID(query string, id int, dest ...interface{}) error {
	err := db.QueryRow(query, id).Scan(dest...)
	if err == sql.ErrNoRows {
		return fmt.Errorf("record, %d, not found", id)
	}
	return err
}

// To ensure the database remains consistent
/* 
    tx, err := s.DB.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    // all your queries using tx instead of s.DB
    return tx.Commit()
*/