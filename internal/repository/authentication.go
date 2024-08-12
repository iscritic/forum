package repository

import (
	"database/sql"
	"forum/internal/entity"
)

func (storage *Storage) CreateUser(user entity.User) error {
	_, err := storage.db.Exec(`INSERT INTO users (username, email, password) VALUES (?, ?, ?)`, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (storage *Storage) GetUserByUsername(username string) (*entity.User, error) {
	row := storage.db.QueryRow(`SELECT id, username, email, password FROM users WHERE username = ?`, username)

	var user entity.User // Создаем переменную user

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password) // Передаем адреса полей структуры для сканирования
	if err != nil {
		return nil, err
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (storage *Storage) GetUserByID(id int) (*entity.User, error) {
	row := storage.db.QueryRow(`SELECT id, username, email, password FROM users WHERE id = ?`, id)

	var user entity.User // Создаем переменную user

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password) // Передаем адреса полей структуры для сканирования
	if err != nil {
		return nil, err
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (storage *Storage) CreateSession(sess entity.Session) error {
	_, err := storage.db.Exec(`INSERT INTO sessions (session_token, user_id, expires_at) VALUES (?, ?, ?)`,
		sess.SessionToken, sess.UserID, sess.ExpiredAt)
	if err != nil {
		return err
	}
	return nil
}

func (storage *Storage) GetSessionByToken(token string) (*entity.Session, error) {
	var session entity.Session
	err := storage.db.QueryRow(`
        SELECT id, session_token, user_id, created_at, expires_at 
        FROM sessions 
        WHERE session_token = ?`, token).Scan(
		&session.ID, &session.SessionToken, &session.UserID, &session.CreatedAt, &session.ExpiredAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No session found
		}
		return nil, err
	}
	return &session, nil
}

func (storage *Storage) DeleteSession(token string) error {
	_, err := storage.db.Exec(`DELETE FROM sessions WHERE session_token = ?`, token)
	return err
}

func (s *Storage) DeleteAllSessionsForUser(userID int) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := s.db.Exec(query, userID)
	return err
}
