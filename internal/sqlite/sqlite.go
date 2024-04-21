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

func (store *Storage) GetAllPosts() ([]*Post, error) {
	// Выполняем запрос к базе данных для выборки всех постов
	rows, err := store.db.Query("SELECT id, title, content, author_id, creation_date FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем слайс для хранения всех постов
	var posts []*Post

	// Итерируем по результатам запроса
	for rows.Next() {
		// Создаем временную переменную для хранения данных поста
		var post Post
		// Сканируем результаты запроса в переменные структуры Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreationDate)
		if err != nil {
			return nil, err
		}
		// Добавляем пост в слайс
		posts = append(posts, &post)
	}

	// Проверяем наличие ошибок после итерации по результатам запроса
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Возвращаем слайс всех постов
	return posts, nil
}

func (store *Storage) GetLastPostID() (int, error) {
	var lastID int

	query := "SELECT id FROM posts ORDER BY id DESC LIMIT 1"
	err := store.db.QueryRow(query).Scan(&lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Если нет строк, возвращаем 0
		}
		log.Println("Error getting last post ID:", err)
		return 0, err
	}

	return lastID, nil
}
