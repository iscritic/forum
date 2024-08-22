package entity

import (
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
