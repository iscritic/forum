package service

import (
	"fmt"
	"forum/internal/entity"
	"forum/internal/repository"
	"strings"
)

func Register(db *repository.Storage, user *entity.User) error {
	err := db.CreateUser(user)
	if err != nil {
		// Check if it's a unique constraint violation for the username or email
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return fmt.Errorf("the username is taken, try another")
		} else if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return fmt.Errorf("the email is taken, try another")
		}
		// For any other error, return the original error
		return err
	}

	return nil
}
