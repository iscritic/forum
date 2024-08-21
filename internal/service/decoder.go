package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"forum/internal/entity"
	"forum/internal/utils"
	"forum/pkg/validator"
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

	if !validator.IsLengthValid(title, 3, 100) ||
		!validator.IsLengthValid(content, 10, 1000) {
		return nil, errors.New("title or content is too long or too short")
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

	if !validator.IsLengthValid(content, 4, 500) {
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

	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if !validator.IsLengthValid(username, 4, 30) {
		return nil, errors.New("username must be between 4 and 30 characters long")
	}

	if !validator.ValidateUsername(username) {
		return nil, errors.New("username can only contain letters, numbers, underscores, and hyphens")
	}

	if !validator.IsLengthValid(email, 1, 256) {
		return nil, errors.New("email is too long")
	}

	if !validator.ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	if !validator.IsLengthValid(password, 8, 50) {
		return nil, errors.New("password must contain at least 8 characters")
	}

	if !validator.ValidatePassword(password) {
		return nil, errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	newUser := &entity.User{
		Username: username,
		Email:    email,
		Password: password,
		Role:     "user",
	}

	return newUser, nil
}
