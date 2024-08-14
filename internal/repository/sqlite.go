package repository

import (
	"database/sql"
	"fmt"
	"os"

	"forum/internal/entity"
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

	querry, err := os.ReadFile("./storage/schema.sql")
	if err != nil {
		return nil, fmt.Errorf("can't read file: %v", err)
	}

	_, err = db.Exec(string(querry))
	if err != nil {
		return nil, fmt.Errorf("can't execute this querry: %v", err)
	}

	return &Storage{db: db}, nil
}

func (storage *Storage) GetAllCategories() ([]entity.Category, error) {
	categories := []entity.Category{}

	rows, err := storage.db.Query("SELECT id, name FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var category entity.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
