package service

import (
	"forum/internal/entity"
	"forum/internal/repository"
)

func GetAllPostRelatedData(db *repository.Storage) ([]entity.PostRelatedData, error) {
	data, err := db.GetAllPostsRelatedData()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetPostRelatedData(db *repository.Storage, id int) (entity.PostRelatedData, error) {
	post, err := db.GetPostRelatedDataByPostID(id)
	if err != nil {
		return entity.PostRelatedData{}, err
	}

	return *post, nil
}

func GetAllLikedPostsById(db *repository.Storage, id int) ([]entity.PostRelatedData, error) {
	data, err := db.GetMyLikedPosts(id)
	if err != nil {
		return nil, err
	}
	// todo fix, internal

	return data, nil
}

func GetCategories(db *repository.Storage) ([]entity.Category, error) {
	categories, err := db.GetAllCategories()
	if err != nil {
		return nil, err
	}
	// todo fix, internal

	return categories, nil
}
