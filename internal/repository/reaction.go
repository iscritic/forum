package repository

import (
	"database/sql"
	"errors"

	"forum/internal/entity"
)

func (storage *Storage) LikePost(userID, postID int) error {
	_, err := storage.db.Exec(`
        INSERT INTO likes (user_id, post_id)
        VALUES (?, ?)`, userID, postID, nil)
	return err
}

func (storage *Storage) LikeComment(userID, commentID int) error {
	_, err := storage.db.Exec(`
        INSERT INTO likes (user_id, comment_id)
        VALUES (?, ?)`, userID, commentID, nil)
	return err
}

func (storage *Storage) DislikePost(userID, postID int) error {
	_, err := storage.db.Exec(`
        INSERT INTO dislikes (user_id, post_id)
        VALUES (?, ?)`, userID, postID, nil)
	return err
}

func (storage *Storage) DislikeComment(userID, commentID int) error {
	_, err := storage.db.Exec(`
        INSERT INTO dislikes (user_id, comment_id)
        VALUES (?, ?)`, userID, commentID, nil)
	return err
}

func (storage *Storage) HasLikedPost(userID, postID int) (bool, error) {
	var count int
	err := storage.db.QueryRow(`
        SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ?`, userID, postID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (storage *Storage) HasDislikedPost(userID, postID int) (bool, error) {
	var count int
	err := storage.db.QueryRow(`
        SELECT COUNT(*) FROM dislikes WHERE user_id = ? AND post_id = ?`, userID, postID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (storage *Storage) RemoveLike(userID, postID int) error {
	_, err := storage.db.Exec(`DELETE FROM likes WHERE user_id = ? AND post_id = ?`, userID, postID)
	return err
}

func (storage *Storage) RemoveDislike(userID, postID int) error {
	_, err := storage.db.Exec(`DELETE FROM dislikes WHERE user_id = ? AND post_id = ?`, userID, postID)
	return err
}

func (storage *Storage) UnlikeComment(userID, commentID int) error {
	_, err := storage.db.Exec(`
        DELETE FROM likes 
        WHERE user_id = ? AND comment_id = ?`, userID, commentID)
	return err
}

func (storage *Storage) UndislikeComment(userID, commentID int) error {
	_, err := storage.db.Exec(`
        DELETE FROM dislikes 
        WHERE user_id = ? AND comment_id = ?`, userID, commentID)
	return err
}

func (storage *Storage) HasLikedComment(userID, commentID int) (bool, error) {
	var count int
	err := storage.db.QueryRow(`
        SELECT COUNT(*) FROM likes WHERE user_id = ? AND comment_id = ?`, userID, commentID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (storage *Storage) HasDislikedComment(userID, commentID int) (bool, error) {
	var count int
	err := storage.db.QueryRow(`
        SELECT COUNT(*) FROM dislikes WHERE user_id = ? AND comment_id = ?`, userID, commentID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (storage *Storage) GetLikesAndDislikesForPost(postID int) (likes, dislikes int, err error) {
	err = storage.db.QueryRow(`
        SELECT 
            (SELECT COUNT(*) FROM likes WHERE post_id = ?) AS likes, 
            (SELECT COUNT(*) FROM dislikes WHERE post_id = ?) AS dislikes`,
		postID, postID).Scan(&likes, &dislikes)
	return
}

func (storage *Storage) GetLikesAndDislikesForComment(commentID int) (likes, dislikes int, err error) {
	err = storage.db.QueryRow(`
        SELECT 
            (SELECT COUNT(*) FROM likes WHERE comment_id = ?) AS likes, 
            (SELECT COUNT(*) FROM dislikes WHERE comment_id = ?) AS dislikes`,
		commentID, commentID).Scan(&likes, &dislikes)
	return
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
