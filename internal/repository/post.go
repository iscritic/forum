package repository

import (
	"database/sql"
	"errors"

	"forum/internal/entity"
)

func (s *Storage) CreatePost(post *entity.Post) (int, error) {
	query := `INSERT INTO posts (title, content, author_id, creation_date) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, post.Title, post.Content, post.AuthorID, post.CreationDate)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, categoryID := range post.CategoryIDs {
		_, err = s.db.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, categoryID)
		if err != nil {
			return 0, err
		}
	}

	return int(postID), nil
}

func (s *Storage) GetAllPost() ([]*entity.Post, error) {
	query := "SELECT id, title, content, author_id, creation_date FROM posts ORDER BY id DESC"

	rows, err := s.db.Query(query)
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

		// Fetch categories for this post
		catRows, err := s.db.Query("SELECT category_id FROM post_categories WHERE post_id = ?", post.ID)
		if err != nil {
			return nil, err
		}
		defer catRows.Close()

		for catRows.Next() {
			var catID int
			if err := catRows.Scan(&catID); err != nil {
				return nil, err
			}
			post.CategoryIDs = append(post.CategoryIDs, catID)
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (Storage *Storage) GetPostByID(id int) (*entity.Post, error) {
	query := "SELECT id, title, content, author_id, category_id, creation_date FROM posts WHERE id = ?"
	row := Storage.db.QueryRow(query, id)

	post := &entity.Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CategoryIDs, &post.CreationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return post, nil
}
