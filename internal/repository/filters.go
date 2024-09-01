package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/internal/entity"
)

func (s Storage) GetAllPostByCategory(categoryID int) ([]*entity.Post, error) {
	query := `
SELECT p.id, p.title, p.content, p.author_id, p.creation_date
FROM posts p
JOIN post_category pc ON p.id = pc.post_id
WHERE pc.category_id = ?;
`

	rows, err := s.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreationDate)
		if err != nil {
			return nil, err
		}
		post.CategoryIDs, err = s.GetCategoryIDsByPostID(post.ID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s Storage) GetAllPostByUser(userID int) ([]*entity.Post, error) {
	query := `
SELECT p.id, p.title, p.content, p.author_id, p.creation_date
FROM posts p
WHERE p.author_id = ?;
`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreationDate)
		if err != nil {
			return nil, err
		}
		post.CategoryIDs, err = s.GetCategoryIDsByPostID(post.ID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s Storage) GetAllLikedPosts(userID int) ([]*entity.Post, error) {
	query := `
	SELECT 
    p.id, p.title, p.content, p.author_id, p.creation_date 
FROM 
    posts p
JOIN
    likes l ON p.id = l.post_id
WHERE 
    l.user_id = ?;
`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreationDate)
		if err != nil {
			return nil, err
		}
		post.CategoryIDs, err = s.GetCategoryIDsByPostID(post.ID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Storage) GetCategoryById(categoryID int) (*entity.Category, error) {
	var category entity.Category

	query := `SELECT id, name FROM category WHERE id = $1`
	err := s.db.QueryRow(query, categoryID).Scan(&category.ID, &category.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &category, nil
}

func (s Storage) GetCategoryIDsByPostID(postID int) ([]int, error) {
	query := "SELECT category_id FROM post_category WHERE post_id = ?"
	rows, err := s.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoryIDs []int
	for rows.Next() {
		var categoryID int
		if err := rows.Scan(&categoryID); err != nil {
			return nil, err
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categoryIDs, nil
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

func (s *Storage) GetLenOfCategories() (int, error) {
	var count int

	row := s.db.QueryRow("SELECT COUNT(*) FROM category")

	err := row.Scan(&count)
	if err != nil {
		fmt.Println("GetLenOfCategories", err)
		return 0, err
	}

	return count, nil
}
