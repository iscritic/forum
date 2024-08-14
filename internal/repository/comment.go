package repository

import "forum/internal/entity"

func (Storage *Storage) CreateComment(comment entity.Comment) error {
	_, err := Storage.db.Exec(`INSERT INTO comments (post_id, content, author_id) VALUES (?, ?, ? )`,
		comment.PostID, comment.Content, comment.AuthorID)
	if err != nil {
		return err
	}

	return nil
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

	// Преобразование [] *Comment в []Comment
	result := make([]entity.Comment, len(comments))
	for i, c := range comments {
		result[i] = *c
	}

	return result, nil
}

//TODO: Edit Comment

//TODO: Delete Comment
