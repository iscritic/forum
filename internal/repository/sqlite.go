package repository

import (
	"database/sql"
	"fmt"
	"os"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", path)
		if err != nil {
			return nil, fmt.Errorf("can't open database: %w", err)
		}
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("can't connect to database: %w", err)
		}
		return &Storage{db: db}, nil
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	querry, err := os.ReadFile("./storage/schema.sql")
	if err != nil {
		return nil, fmt.Errorf("can't read file: %v", err)
	}
	migration, err := os.ReadFile("./storage/migration.sql")
	if err != nil {
		return nil, fmt.Errorf("can't read file: %v", err)
	}

	_, err = db.Exec(string(querry))
	if err != nil {
		return nil, fmt.Errorf("can't execute this querry: %v", err)
	}
	_, err = db.Exec(string(migration))
	if err != nil {
		return nil, fmt.Errorf("can't execute this querry: %v", err)
	}

	return &Storage{db: db}, nil
}
