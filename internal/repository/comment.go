package repository

import (
	"database/sql"
	"errors"
	"forum/internal/entity"
)

func (Storage *Storage) CreateComment(comment entity.Comment) error {
	_, err := Storage.db.Exec(`INSERT INTO comments (post_id, content, author_id) VALUES (?, ?, ? )`,
		comment.PostID, comment.Content, comment.AuthorID)
	if err != nil {
		return err
	}

	return nil
}

func (Storage *Storage) GetCommentByID(id int) (*entity.Comment, error) {
	query := "SELECT id, post_id, content, author_id, creation_date FROM comments WHERE id = ?"
	row := Storage.db.QueryRow(query, id)

	comment := &entity.Comment{}
	err := row.Scan(&comment.ID, &comment.PostID, &comment.Content, &comment.AuthorID, &comment.CreationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return comment, nil
}

func (s *Storage) GetAllComments(postID int) ([]*entity.Comment, error) {
	rows, err := s.db.Query("SELECT id, content, creation_date FROM comments WHERE post_id = ? ORDER BY id DESC", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.Comment

	for rows.Next() {
		var comment entity.Comment
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

func (s *Storage) GetCommentsByPostID(postID int) ([]entity.Comment, error) {
	rows, err := s.db.Query(`
       SELECT id, post_id, content, author_id, likes, dislikes, creation_date
       FROM comments
       WHERE post_id = ?
       ORDER BY creation_date DESC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.Comment

	for rows.Next() {
		var comment entity.Comment
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

	result := make([]entity.Comment, len(comments))
	for i, c := range comments {
		result[i] = *c
	}

	return result, nil
}

func (s *Storage) UpdateComment(comment *entity.Comment) error {
	_, err := s.db.Exec(`UPDATE comments SET content = ?, creation_date = CURRENT_TIMESTAMP WHERE id = ? AND author_id = ?`,
		comment.Content, comment.ID, comment.AuthorID)
	return err
}

func (s *Storage) DeleteComment(commentID int) error {
	_, err := s.db.Exec("DELETE FROM comments WHERE id = ?", commentID)
	return err
}
