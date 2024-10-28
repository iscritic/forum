package repository

import (
	"fmt"
	"forum/internal/entity"
)

func (s *Storage) CreateNotification(userID int, message string, postID *int, commentID *int) error {
	_, err := s.db.Exec(`INSERT INTO notifications (user_id, message, post_id, comment_id) VALUES (?, ?, ?, ?)`, userID, message, postID, commentID)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

func (s *Storage) GetUnreadNotifications(userID int) ([]*entity.Notification, error) {
	rows, err := s.db.Query(`SELECT id, message, post_id, comment_id, created_at FROM notifications WHERE user_id = ? AND read = false`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*entity.Notification
	for rows.Next() {
		var n entity.Notification
		if err := rows.Scan(&n.ID, &n.Message, &n.PostID, &n.CommentID, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	return notifications, nil
}

func (s *Storage) MarkNotificationAsRead(notificationID int) error {
	_, err := s.db.Exec(`UPDATE notifications SET read = 1 WHERE id = ?`, notificationID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}
