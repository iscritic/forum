package repository

import (
	"database/sql"
	"fmt"
)

func (s *Storage) GetPostAuthor(postID int) (int, error) {
	var authorID int
	query := `SELECT author_id FROM posts WHERE id = ?`
	err := s.db.QueryRow(query, postID).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no post found with ID %d", postID)
		}
		return 0, fmt.Errorf("failed to get author for post ID %d: %w", postID, err)
	}
	return authorID, nil
}
