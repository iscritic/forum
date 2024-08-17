package service

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"forum/internal/utils"
	"forum/internal/utils/validator"
	"net/http"
	"reflect"
)

// DecodePost decodes and validates post data from an HTTP request form
func DecodePost(r *http.Request) (*entity.Post, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("error parsing form: %w", err)
	}

	title := r.Form.Get("title")
	content := r.Form.Get("content")
	categoryIDStr := r.Form.Get("category")

	categoryID, err := utils.Etoi(categoryIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	// Basic validation using your custom validator
	if validator.IsEmptyValue(reflect.ValueOf(title)) ||
		validator.IsEmptyValue(reflect.ValueOf(content)) ||
		validator.IsEmptyValue(reflect.ValueOf(categoryIDStr)) {
		return nil, fmt.Errorf("title, content, and category are required")
	}

	if !validator.IsLengthValid(title, 0, 100) ||
		!validator.IsLengthValid(content, 0, 1000) {
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
func DecodeComment(r *http.Request) (*entity.Post, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("error parsing form: %w", err)
	}

	postIDStr := r.Form.Get("post_id")
	postID, err := utils.Etoi(postIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %w", err)
	}

	//
	//comment.PostID = postID
	//comment.Content = r.Form.Get("content")
	//
	//comment.AuthorID = r.Context().Value("userID").(int)
	//
	//// TODO author comment id

	//title := r.Form.Get("title")
	//content := r.Form.Get("content")
	//categoryIDStr := r.Form.Get("category")
	//
	//categoryID, err := utils.Etoi(categoryIDStr)
	//if err != nil {
	//	return nil, fmt.Errorf("invalid category ID: %w", err)
	//}
	//
	//// Basic validation using your custom validator
	//if validator.IsEmptyValue(reflect.ValueOf(title)) ||
	//	validator.IsEmptyValue(reflect.ValueOf(content)) ||
	//	validator.IsEmptyValue(reflect.ValueOf(categoryIDStr)) {
	//	return nil, fmt.Errorf("title, content, and category are required")
	//}
	//
	//if !validator.IsLengthValid(title, 0, 100) ||
	//	!validator.IsLengthValid(content, 0, 1000) {
	//	return nil, errors.New("title or content is too long")
	//}
	//
	//post := &entity.Post{
	//	Title:      title,
	//	Content:    content,
	//	CategoryID: categoryID,
	//	AuthorID:   r.Context().Value("userID").(int),
	//}

	return post, nil
}
