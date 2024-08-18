package service

import (
	"context"
	"fmt"

	"forum/internal/entity"
	"forum/internal/repository"
)

type HomePageData struct {
	Posts      []*entity.PostRelatedData
	Categories []entity.Category
	UserInfo   *entity.User
}

func getUserInfo(ctx context.Context, storage *repository.Storage) (*entity.User, error) {
	role, ok := ctx.Value("role").(string)
	if !ok {
		return nil, fmt.Errorf("role is not a string")
	}

	if role != "guest" {
		userID, ok := ctx.Value("userID").(int)
		if !ok {
			return nil, fmt.Errorf("userID is not an int")
		}

		userInfo, err := storage.GetUserByID(userID)
		if err != nil {
			return nil, err
		}
		return userInfo, nil
	}

	return &entity.User{Role: role}, nil
}

func GetHomePageData(db *repository.Storage, ctx context.Context) (HomePageData, error) {
	posts, err := GetAllPostRelatedData(db)
	if err != nil {
		return HomePageData{}, err
	}

	categories, err := GetCategories(db)
	if err != nil {
		return HomePageData{}, err
	}

	userInfo, err := getUserInfo(ctx, db)
	if err != nil {
		return HomePageData{}, err
	}

	return HomePageData{
		Posts:      posts,
		Categories: categories,
		UserInfo:   userInfo,
	}, nil
}
