package repository

import (
	"database/sql"
	"errors"

	"forum/internal/entity"
)

func (Storage *Storage) CreatePost(post *entity.Post) (int, error) {
	res, err := Storage.db.Exec(`INSERT INTO posts (title, content, author_id) VALUES (?, ?, ?)`,
		post.Title, post.Content, post.AuthorID)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	for _, categoryID := range post.CategoryIDs {
		_, err := Storage.db.Exec(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, lastID, categoryID)
		if err != nil {
			return 0, err
		}
	}
	return int(lastID), nil
}

func (Storage *Storage) GetAllPost() ([]*entity.Post, error) {
	query := "SELECT id, title, content, author_id, creation_date FROM posts ORDER BY id DESC"

	rows, err := Storage.db.Query(query)
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

		categoryQuery := "SELECT category_id FROM post_category WHERE post_id = ?"
		categoryRows, err := Storage.db.Query(categoryQuery, post.ID)
		if err != nil {
			return nil, err
		}
		defer categoryRows.Close()

		var categoryIDs []int
		for categoryRows.Next() {
			var categoryID int
			err := categoryRows.Scan(&categoryID)
			if err != nil {
				return nil, err
			}
			categoryIDs = append(categoryIDs, categoryID)
		}
		post.CategoryIDs = categoryIDs
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (Storage *Storage) GetPostByID(id int) (*entity.Post, error) {
	query := "SELECT id, title, content, author_id, creation_date FROM posts WHERE id = ?"
	row := Storage.db.QueryRow(query, id)

	post := &entity.Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	categoryQuery := "SELECT category_id FROM post_category WHERE post_id = ?"
	categoryRows, err := Storage.db.Query(categoryQuery, post.ID)
	if err != nil {
		return nil, err
	}
	defer categoryRows.Close()

	var categoryIDs []int
	for categoryRows.Next() {
		var categoryID int
		err := categoryRows.Scan(&categoryID)
		if err != nil {
			return nil, err
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	post.CategoryIDs = categoryIDs

	return post, nil
}

func (Storage *Storage) UpdatePost(post *entity.Post) error {
	_, err := Storage.db.Exec(`UPDATE posts SET title = ?, content = ? WHERE id = ? AND author_id = ?`,
		post.Title, post.Content, post.ID, post.AuthorID)
	if err != nil {
		return err
	}

	_, err = Storage.db.Exec(`DELETE FROM post_category WHERE post_id = ?`, post.ID)
	if err != nil {
		return err
	}

	for _, categoryID := range post.CategoryIDs {
		_, err := Storage.db.Exec(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`, post.ID, categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) DeletePost(postID int) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE id = ?", postID)
	return err
}
