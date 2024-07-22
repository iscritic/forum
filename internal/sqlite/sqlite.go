package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
)

type Post struct {
	ID           int
	Title        string
	Content      string
	AuthorID     int
	Category     string
	CreationDate time.Time
}

type Comment struct {
	ID           int
	PostID       int
	Content      string
	AuthorID     int
	CreationDate time.Time
}

type PostData struct {
	Post    Post
	Comment []*Comment
}

type User struct {
	ID           int
	Username     string
	Email        string
	Password     string
	Role         string
	CreationDate time.Time
}

type Category struct {
	ID   int
	Name string
}

type Session struct {
	ID           int
	SessionToken uuid.UUID
	UserID       int
	CreatedAt    time.Time
	ExpiredAt    time.Time
}

type Storage struct {
	db *sql.DB
}

var queries = []string{
	`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY,
        username TEXT UNIQUE,
        email TEXT UNIQUE,
        password TEXT,
        role TEXT,
        creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`,
	`CREATE TABLE IF NOT EXISTS category (
        id INTEGER PRIMARY KEY,
        name TEXT UNIQUE
    )`,
	`CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY,
        title TEXT,
        content TEXT,
        author_id INTEGER,
        category TEXT,
        creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (author_id) REFERENCES users(id),
        FOREIGN KEY (category) REFERENCES category(name)
    )`,
	`CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY,
		post_id INTEGER,
		content TEXT,
		author_id INTEGER,
		creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts(id),
        FOREIGN KEY (author_id) REFERENCES users(id)
	)`,
	`CREATE TABLE IF NOT EXISTS category (
		id INTEGER PRIMARY KEY,
		name TEXT UNIQUE
	)`,

	`CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_token TEXT UNIQUE NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
)`,
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	for _, querie := range queries {
		_, err := db.Exec(querie)
		if err != nil {
			log.Fatalf("can't execute queries %q: %w", querie, err)
		}
	}

	return &Storage{db: db}, nil
}

func (Storage *Storage) CreatePost(post Post) (int, error) {
	res, err := Storage.db.Exec(`INSERT INTO posts (title, content) VALUES (?, ?)`,
		post.Title, post.Content)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastID), nil
}

func (Storage *Storage) GetPostByID(id int) (*Post, error) {
	// Предполагается, что у вас есть поле DB типа *sql.DB в вашей структуре Application
	query := "SELECT id, title, content,  creation_date FROM posts WHERE id = ?"
	row := Storage.db.QueryRow(query, id)

	post := &Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CreationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Пост не найден
		}
		return nil, err // Произошла ошибка при выполнении запроса
	}

	return post, nil
}

func (Storage *Storage) GetAllPosts() ([]*Post, error) {
	// TODO authors ids

	rows, err := Storage.db.Query("SELECT  id, title, content, creation_date FROM posts ORDER BY id DESC ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем слайс для хранения всех постов
	var posts []*Post

	// Итерируем по результатам запроса
	for rows.Next() {
		// Создаем новую переменную для хранения данных поста на каждой итерации
		var post Post
		// Сканируем результаты запроса в переменные структуры Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreationDate)
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

func (Storage *Storage) CreateComment(comment Comment) error {
	_, err := Storage.db.Exec(`INSERT INTO comments (post_id, content) VALUES (?, ? )`,
		comment.PostID, comment.Content)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetAllComments(postID int) ([]*Comment, error) {
	rows, err := s.db.Query("SELECT id, content, creation_date FROM comments WHERE post_id = ? ORDER BY id DESC", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment

	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreationDate); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (storage *Storage) CreateUser(user User) error {
	_, err := storage.db.Exec(`INSERT INTO users (username, email, password) VALUES (?, ?, ?)`, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (storage *Storage) GetUserByUsername(username string) (*User, error) {
	row := storage.db.QueryRow(`SELECT id, username, email, password FROM users WHERE username = ?`, username)

	var user User // Создаем переменную user

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password) // Передаем адреса полей структуры для сканирования
	if err != nil {
		return nil, err
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (storage *Storage) GetUserByID(id int) (*User, error) {
	row := storage.db.QueryRow(`SELECT id, username, email, password FROM users WHERE id = ?`, id)

	var user User // Создаем переменную user

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password) // Передаем адреса полей структуры для сканирования
	if err != nil {
		return nil, err
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (storage *Storage) CreateSession(sess Session) error {
	_, err := storage.db.Exec(`INSERT INTO sessions (session_token, user_id, expires_at) VALUES (?, ?, ?)`,
		sess.SessionToken, sess.UserID, sess.ExpiredAt)
	if err != nil {
		return err
	}
	return nil
}

func (storage *Storage) GetSessionByToken(token string) (*Session, error) {
	var session Session
	err := storage.db.QueryRow(`
        SELECT id, session_token, user_id, created_at, expires_at 
        FROM sessions 
        WHERE session_token = ?`, token).Scan(
		&session.ID, &session.SessionToken, &session.UserID, &session.CreatedAt, &session.ExpiredAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No session found
		}
		return nil, err
	}
	return &session, nil
}

func (storage *Storage) DeleteSession(token string) error {
	_, err := storage.db.Exec(`DELETE FROM sessions WHERE session_token = ?`, token)
	return err
}
