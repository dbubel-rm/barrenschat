// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package db

import (
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	Database *sqlx.DB
}

func New(file string) (*SQLite, error) {
	sqlite, err := sqlx.Connect("sqlite3", file)

	if err != nil {
		return nil, errors.New("Error connecting to DB when application started")
	}

	// Arbitrary but needed
	sqlite.SetMaxIdleConns(10)
	sqlite.DB.SetMaxIdleConns(10)
	sqlite.DB.SetMaxOpenConns(10)

	db := SQLite{
		Database: sqlite,
	}

	return &db, nil
}

// Close closes a DB value being used with MongoDB.
func (db *SQLite) Close() {
	db.Database.Close()
}
