package service

import (
	"forum/internal/entity"
	"forum/internal/repository"
)

func GetAllPostRelatedDataByCategory(db *repository.Storage, categoryID int) ([]*entity.PostRelatedData, error) {
	data, err := db.GetAllPostRelatedDataByCategory(categoryID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetAllPostRelatedDataByUserID(db *repository.Storage, userID int) ([]*entity.PostRelatedData, error) {
	data, err := db.GetAllPostRelatedDataByUser(userID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetMyLikedPosts(db *repository.Storage, userID int) ([]*entity.PostRelatedData, error) {
	data, err := db.GetMyLikedPosts(userID)
	if err != nil {
		return nil, err
	}

	return data, nil
}
