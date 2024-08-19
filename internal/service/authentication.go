package service

import (
	"fmt"
	"strings"

	"forum/internal/entity"
	"forum/internal/repository"
)

func Register(db *repository.Storage, user *entity.User) error {
	err := db.CreateUser(user)
	if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
		return fmt.Errorf("the username is taken, try another")
	} else if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
		return fmt.Errorf("the email is taken, try another")
	} else if err != nil {
		return err
	}

	return nil
}
