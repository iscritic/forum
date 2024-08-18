package service

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"forum/internal/utils"
	"forum/pkg/validator"
	"net/http"
	"strings"
)

// DecodePost decodes and validates post data from an HTTP request form
func DecodePost(r *http.Request) (*entity.Post, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("error parsing form: %w", err)
	}

	title := strings.TrimSpace(r.Form.Get("title"))
	content := strings.TrimSpace(r.Form.Get("content"))
	categoryIDStr := r.Form.Get("category")

	categoryID, err := utils.Etoi(categoryIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	if !validator.IsLengthValid(title, 1, 100) ||
		!validator.IsLengthValid(content, 1, 1000) {
		return nil, errors.New("title or content is too long")
	}

	post := &entity.Post{
		Title:      title,
		Content:    content,
		CategoryID: categoryID,
		AuthorID:   r.Context().Value("userID").(int),
	}

	return post, nil
}

// DecodeComment decodes and validates comment data from an HTTP request form
func DecodeComment(r *http.Request) (*entity.Comment, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("error parsing form: %w", err)
	}

	postIDStr := r.Form.Get("post_id")
	postID, err := utils.Etoi(postIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}

	content := strings.TrimSpace(r.Form.Get("content"))

	if !validator.IsLengthValid(content, 0, 500) {
		return nil, errors.New("content is too long")
	}

	comment := &entity.Comment{
		PostID:   postID,
		Content:  content,
		AuthorID: r.Context().Value("userID").(int),
	}

	return comment, nil

}

func DecodeUser(r *http.Request) (*entity.User, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("error parsing form: %w", err)
	}

	username := strings.TrimSpace(r.Form.Get("username"))
	email := strings.TrimSpace(r.Form.Get("email"))
	password := strings.TrimSpace(r.Form.Get("password"))

	if !validator.ValidateUsername(username) ||
		!validator.ValidatePassword(password) {
		return nil, errors.New("invalid username or password")
	}

	if !validator.ValidateEmail(email) {
		return nil, fmt.Errorf("email is invalid")
	}

	if !validator.IsLengthValid(username, 4, 30) ||
		!validator.IsLengthValid(password, 8, 50) ||
		!validator.IsLengthValid(email, 1, 256) {
		return nil, errors.New("username, password, or email is too long")
	}

	newUser := &entity.User{
		Username: username,
		Email:    email,
		Password: password,
		Role:     "user",
	}

	return newUser, nil
}
