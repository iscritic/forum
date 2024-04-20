package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Post struct {
	ID           int
	Title        string
	Content      string
	AuthorID     int
	CreationDate time.Time
}

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

func (store *Storage) CreatePost(post Post) error {
	_, err := store.db.Exec(`INSERT INTO posts (title, content, author_id) VALUES (?, ?, ?)`,
		post.Title, post.Content, post.AuthorID)
	if err != nil {
		return err
	}
	return nil
}

func (store *Storage) GetPostByID(id int) (*Post, error) {
	// Предполагается, что у вас есть поле DB типа *sql.DB в вашей структуре Application
	query := "SELECT id, title, content, author_id, creation_date FROM posts WHERE id = ?"
	row := store.db.QueryRow(query, id)

	post := &Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreationDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пост не найден
		}
		return nil, err // Произошла ошибка при выполнении запроса
	}

	return post, nil
}
