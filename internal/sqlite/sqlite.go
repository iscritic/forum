package sqlite

import (
	"database/sql"
	"fmt"
	"log"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	// Создаем таблицу пользователей.
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
                        id INTEGER PRIMARY KEY,
                        username TEXT UNIQUE,
                        email TEXT UNIQUE,
                        password TEXT,
                        creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                    )`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS posts (
                        id INTEGER PRIMARY KEY,
                        title TEXT,
                        content TEXT,
                        author_id INTEGER,
                        creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                    )`)
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{db: db}, nil

}
