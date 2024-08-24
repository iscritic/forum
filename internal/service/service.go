package service

import (
	"context"

	"forum/internal/entity"
	"forum/internal/repository"
)

func GetAllPostRelatedData(db *repository.Storage) ([]*entity.PostRelatedData, error) {
	data, err := db.GetAllPostsRelatedData()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetPostRelatedData(ctx context.Context, db *repository.Storage, id int) (*entity.PostRelatedData, error) {
	post, err := db.GetPostRelatedDataByPostID(id)
	if err != nil {
		return nil, err
	}

	// get user id from context for checking login
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return post, err
	}
	ctxUser, err := db.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	post.CtxUser = *ctxUser

	return post, nil
}

func GetCategories(db *repository.Storage) ([]entity.Category, error) {
	categories, err := db.GetAllCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}
