package service

import (
	"context"

	"forum/internal/entity"
	"forum/internal/repository"
)

type dataPage struct {
	Post       entity.Post
	Categories []entity.Category
	UserInfo   *entity.User
	Comments   []entity.CommentRelatedData
}

func GetAllPostRelatedData(db *repository.Storage) ([]*entity.PostRelatedData, error) {
	data, err := db.GetAllPostsRelatedData()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetPostRelatedData(ctx context.Context, db *repository.Storage, id int) (dataPage, error) {
	post, err := db.GetPostByID(id)
	if err != nil {
		return dataPage{}, err
	}

	categories, err := GetCategories(db)
	if err != nil {
		return dataPage{}, err
	}
	user, err := GetUserInfo(ctx, db)
	if err != nil {
		return dataPage{}, err
	}
	comments, err := db.GetAllCommentsRelatedDataByPostID(id)
	if err != nil {
		return dataPage{}, err
	}
	data := dataPage{
		Post:       *post,
		Categories: categories,
		UserInfo:   user,
		Comments:   comments,
	}

	return data, nil
}

func GetCategories(db *repository.Storage) ([]entity.Category, error) {
	categories, err := db.GetAllCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}
