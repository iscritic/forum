package repository

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

func (Storage *Storage) CreatePost(post entity.Post) (int, error) {
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

func (Storage *Storage) GetAllPost() ([]*entity.Post, error) {

	query := "SELECT  id, title, content, author_id, category_id, creation_date FROM posts ORDER BY id DESC"

	rows, err := Storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CategoryID, &post.CreationDate)
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

func (Storage *Storage) GetPostByID(id int) (*entity.Post, error) {
	query := "SELECT id, title, content, author_id, category_id, creation_date FROM posts WHERE id = ?"
	row := Storage.db.QueryRow(query, id)

	post := &entity.Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CategoryID, &post.CreationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return post, nil
}

//TODO: Edit Post

//TODO: Delete Post
