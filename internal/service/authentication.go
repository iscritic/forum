package service

import (
	"forum/internal/entity"
	"forum/internal/repository"
)

func Register(db *repository.Storage, user entity.User) error {

	//TODO: add validation

	//TODO: checking duplicates

	err := db.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}
