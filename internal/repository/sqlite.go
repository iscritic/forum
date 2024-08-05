package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gofrs/uuid"
)

type Post struct {
	ID           int
	Title        string
	Content      string
	AuthorID     int
	CategoryID   int
	Likes        int
	Dislikes     int
	CreationDate time.Time
}

type Comment struct {
	ID           int
	PostID       int
	Content      string
	AuthorID     int
	Likes        int
	Dislikes     int
	CreationDate time.Time
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

type Like struct {
	ID        int
	PostID    int
	CommentID int
	UserID    int
	Grade     int
}

type PostRelatedData struct {
	Post     Post
	CommentR []CommentRelatedData
	User     User
	Category Category
}

type CommentRelatedData struct {
	Comment Comment
	User    User
}

type PostData struct {
	Post     PostRelatedData
	Comments []*CommentRelatedData
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

	querry, err := ioutil.ReadFile("./storage/schema.sql")
	if err != nil {
		return nil, fmt.Errorf("can't read file: %v", err)
	}

	_, err = db.Exec(string(querry))
	if err != nil {
		return nil, fmt.Errorf("can't execute this querry: %v", err)
	}

	return &Storage{db: db}, nil
}

func (Storage *Storage) CreatePost(post Post) (int, error) {
	res, err := Storage.db.Exec(`INSERT INTO posts (title, content, author_id, category_id) VALUES (?, ?, ?, ?)`,
		post.Title, post.Content, post.AuthorID, post.CategoryID)
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
	query := "SELECT id, title, content, author_id, category_id, creation_date FROM posts WHERE id = ?"
	row := Storage.db.QueryRow(query, id)

	post := &Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CategoryID, &post.CreationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Пост не найден
		}
		return nil, err // Произошла ошибка при выполнении запроса
	}

	return post, nil
}

func (s *Storage) GetPostRelatedData(postID int) (*PostRelatedData, error) {
	// Получение поста по ID
	post, err := s.GetPostByID(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	// Получение связанных комментариев
	comments, err := s.GetCommentsRelatedData(postID)
	if err != nil {
		return nil, err
	}

	// Получение информации об авторе поста
	author, err := s.GetUserByID(post.AuthorID)
	if err != nil {
		return nil, err
	}

	// Получение информации о категории
	categoryName, err := s.GetCategoryById(post.CategoryID)
	if err != nil {
		return nil, err
	}

	// Заполнение и возврат структуры PostRelatedData
	postRelatedData := &PostRelatedData{
		Post:     *post,
		CommentR: comments,
		User:     *author,
		Category: Category{ID: post.CategoryID, Name: categoryName},
	}

	return postRelatedData, nil
}

func (Storage *Storage) GetAllPosts() ([]*Post, error) {
	// TODO authors ids

	rows, err := Storage.db.Query("SELECT  id, title, content, author_id, category_id, creation_date FROM posts ORDER BY id DESC ")
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
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CategoryID, &post.CreationDate)
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
	_, err := Storage.db.Exec(`INSERT INTO comments (post_id, content, author_id) VALUES (?, ?, ? )`,
		comment.PostID, comment.Content, comment.AuthorID)
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

func (s *Storage) GetCommentsRelatedData(postID int) ([]CommentRelatedData, error) {
	// Получаем все комментарии
	rows, err := s.db.Query(`
        SELECT id, post_id, content, author_id, likes, dislikes, creation_date 
        FROM comments 
        WHERE post_id = ? 
        ORDER BY creation_date DESC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Content,
			&comment.AuthorID,
			&comment.Likes,
			&comment.Dislikes,
			&comment.CreationDate,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Получаем уникальные ID авторов
	authorIDs := make(map[int]struct{})
	for _, comment := range comments {
		authorIDs[comment.AuthorID] = struct{}{}
	}

	// Получаем информацию о пользователях
	var authors []User
	for authorID := range authorIDs {
		author, err := s.GetUserByID(authorID)
		if err != nil {
			return nil, err
		}
		authors = append(authors, *author)
	}

	// Создаем отображение ID авторов для быстрого доступа
	authorMap := make(map[int]User)
	for _, author := range authors {
		authorMap[author.ID] = author
	}

	// Создаем результат
	var commentsRelatedData []CommentRelatedData
	for _, comment := range comments {
		user := authorMap[comment.AuthorID]
		commentsRelatedData = append(commentsRelatedData, CommentRelatedData{
			Comment: comment,
			User:    user,
		})
	}

	return commentsRelatedData, nil
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

func (storage *Storage) GetCountOfLikes(postID int) (int, error) {
	var count int

	return count, nil
}

func (storage *Storage) GetCountOfDislikes(postID int) (int, error) {
	var count int

	return count, nil
}

func (storage *Storage) GetCategoryById(categoryID int) (string, error) {
	var category string

	return category, nil
}

func (s *Storage) GetPostRelatedDataByID(postID int) (*PostRelatedData, error) {
	// 1. Fetch the post by ID
	post, err := s.GetPostByID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post: %w", err) // Wrap the error
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	// 2. Fetch related comments
	comments, err := s.GetCommentsRelatedData(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err) // Wrap the error
	}

	// 3. Fetch author information
	author, err := s.GetUserByID(post.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch author: %w", err) // Wrap the error
	}

	// 4. Fetch category information
	categoryName, err := s.GetCategoryById(post.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category: %w", err) // Wrap the error
	}

	// 5. Construct and return the PostRelatedData struct
	postRelatedData := &PostRelatedData{
		Post:     *post,
		CommentR: comments,
		User:     *author,
		Category: Category{ID: post.CategoryID, Name: categoryName},
	}

	return postRelatedData, nil
}

func (s *Storage) GetCommentsByPostID(postID int) ([]Comment, error) {
	rows, err := s.db.Query(`
        SELECT id, post_id, content, author_id, likes, dislikes, creation_date 
        FROM comments 
        WHERE post_id = ? 
        ORDER BY creation_date DESC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment

	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Content,
			&comment.AuthorID,
			&comment.Likes,
			&comment.Dislikes,
			&comment.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Преобразование [] *Comment в []Comment
	result := make([]Comment, len(comments))
	for i, c := range comments {
		result[i] = *c
	}

	return result, nil
}

func (s *Storage) GetPostsRelatedData() ([]PostRelatedData, error) {
	// Получение всех постов
	posts, err := s.GetAllPosts()
	if err != nil {
		return nil, err
	}

	// Слайс для хранения всех связанных данных
	var postsRelatedData []PostRelatedData

	// Итерация по каждому посту для получения связанных данных
	for _, post := range posts {
		// Получение связанных комментариев
		comments, err := s.GetCommentsRelatedData(post.ID)
		if err != nil {
			return nil, err
		}

		// Получение информации об авторе поста
		author, err := s.GetUserByID(post.AuthorID)
		if err != nil {
			return nil, err
		}

		// Получение информации о категории
		categoryName, err := s.GetCategoryById(post.CategoryID)
		if err != nil {
			return nil, err
		}

		// Заполнение структуры PostRelatedData для текущего поста
		postRelatedData := PostRelatedData{
			Post:     *post,
			CommentR: comments,
			User:     *author,
			Category: Category{ID: post.CategoryID, Name: categoryName},
		}

		// Добавление в слайс
		postsRelatedData = append(postsRelatedData, postRelatedData)
	}

	return postsRelatedData, nil
}
