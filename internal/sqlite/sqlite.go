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
	Likes        int
	Dislikes     int
}

type Comment struct {
	ID           int
	PostID       int
	Content      string
	AuthorID     int
	CreationDate time.Time
	Likes        int
	Dislikes     int
}

type PostData struct {
	Post     Post
	Comment  []*Comment
	Likes    int
	Dislikes int
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
	`CREATE TABLE IF NOT EXISTS likes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	post_id INTEGER,
	comment_id INTEGER,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (post_id) REFERENCES posts(id),
	FOREIGN KEY (comment_id) REFERENCES comments(id),
	UNIQUE(user_id, post_id, comment_id)
)`,
	`CREATE TABLE IF NOT EXISTS dislikes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	post_id INTEGER,
	comment_id INTEGER,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (post_id) REFERENCES posts(id),
	FOREIGN KEY (comment_id) REFERENCES comments(id),
	UNIQUE(user_id, post_id, comment_id)
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

func (storage *Storage) GetPostByID(id int) (*Post, error) {
	query := `
        SELECT 
            p.id, p.title, p.content, p.creation_date,
            (SELECT COUNT(*) FROM likes WHERE post_id = p.id) AS likes,
            (SELECT COUNT(*) FROM dislikes WHERE post_id = p.id) AS dislikes
        FROM posts p
        WHERE p.id = ?`
	row := storage.db.QueryRow(query, id)

	post := &Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CreationDate, &post.Likes, &post.Dislikes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Post not found
		}
		return nil, err // Error occurred
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
	rows, err := s.db.Query(`
        SELECT 
            c.id, c.content, c.creation_date,
            (SELECT COUNT(*) FROM likes WHERE comment_id = c.id) AS likes,
            (SELECT COUNT(*) FROM dislikes WHERE comment_id = c.id) AS dislikes
        FROM comments c
        WHERE c.post_id = ? 
        ORDER BY c.id DESC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreationDate, &comment.Likes, &comment.Dislikes); err != nil {
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

func (storage *Storage) LikePost(userID, postID int) error {
	_, err := storage.db.Exec(`
        INSERT INTO likes (user_id, post_id)
        VALUES (?, ?)
        ON CONFLICT(user_id, post_id, comment_id) DO NOTHING`, userID, postID, nil)
	return err
}

func (storage *Storage) LikeComment(userID, commentID int) error {
	_, err := storage.db.Exec(`
        INSERT INTO likes (user_id, comment_id)
        VALUES (?, ?)
        ON CONFLICT(user_id, post_id, comment_id) DO NOTHING`, userID, commentID, nil)
	return err
}

func (storage *Storage) DislikePost(userID, postID int) error {
	_, err := storage.db.Exec(`
        INSERT INTO dislikes (user_id, post_id)
        VALUES (?, ?)
        ON CONFLICT(user_id, post_id, comment_id) DO NOTHING`, userID, postID, nil)
	return err
}

func (storage *Storage) DislikeComment(userID, commentID int) error {
	_, err := storage.db.Exec(`
        INSERT INTO dislikes (user_id, comment_id)
        VALUES (?, ?)
        ON CONFLICT(user_id, post_id, comment_id) DO NOTHING`, userID, commentID, nil)
	return err
}

func (storage *Storage) HasLikedPost(userID, postID int) (bool, error) {
	var count int
	err := storage.db.QueryRow(`
        SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ?`, userID, postID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (storage *Storage) HasDislikedPost(userID, postID int) (bool, error) {
	var count int
	err := storage.db.QueryRow(`
        SELECT COUNT(*) FROM dislikes WHERE user_id = ? AND post_id = ?`, userID, postID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (storage *Storage) RemoveLike(userID, postID int) error {
	_, err := storage.db.Exec(`DELETE FROM likes WHERE user_id = ? AND post_id = ?`, userID, postID)
	return err
}

func (storage *Storage) RemoveDislike(userID, postID int) error {
	_, err := storage.db.Exec(`DELETE FROM dislikes WHERE user_id = ? AND post_id = ?`, userID, postID)
	return err
}

func (storage *Storage) UnlikeComment(userID, commentID int) error {
	_, err := storage.db.Exec(`
        DELETE FROM likes 
        WHERE user_id = ? AND comment_id = ?`, userID, commentID)
	return err
}

func (storage *Storage) UndislikeComment(userID, commentID int) error {
	_, err := storage.db.Exec(`
        DELETE FROM dislikes 
        WHERE user_id = ? AND comment_id = ?`, userID, commentID)
	return err
}

func (storage *Storage) GetLikesAndDislikesForPost(postID int) (likes, dislikes int, err error) {
	err = storage.db.QueryRow(`
        SELECT 
            (SELECT COUNT(*) FROM likes WHERE post_id = ?) AS likes, 
            (SELECT COUNT(*) FROM dislikes WHERE post_id = ?) AS dislikes`,
		postID, postID).Scan(&likes, &dislikes)
	return
}

func (storage *Storage) GetLikesAndDislikesForComment(commentID int) (likes, dislikes int, err error) {
	err = storage.db.QueryRow(`
        SELECT 
            (SELECT COUNT(*) FROM likes WHERE comment_id = ?) AS likes, 
            (SELECT COUNT(*) FROM dislikes WHERE comment_id = ?) AS dislikes`,
		commentID, commentID).Scan(&likes, &dislikes)
	return
}
