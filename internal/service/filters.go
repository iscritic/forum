package service

import (
	"context"

	"forum/internal/entity"
	"forum/internal/repository"
)

type PageData struct {
	Posts      []*entity.PostRelatedData
	Categories []entity.Category
	UserInfo   entity.User
}

func GetAllPostRelatedDataByCategory(ctx context.Context, db *repository.Storage, categoryID int) (PageData, error) {
	posts, err := db.GetAllPostRelatedDataByCategory(categoryID)
	if err != nil {
		return PageData{}, err
	}

	categories, err := GetCategories(db)
	if err != nil {
		return PageData{}, err
	}

	userInfo, err := GetUserInfo(ctx, db)
	if err != nil {
		return PageData{}, err
	}

	return PageData{
		Posts:      posts,
		Categories: categories,
		UserInfo:   *userInfo,
	}, nil
}

func GetMyLikedPosts(ctx context.Context, db *repository.Storage, userID int) (PageData, error) {
	posts, err := db.GetMyLikedPosts(userID)
	if err != nil {
		return PageData{}, err
	}

	categories, err := GetCategories(db)
	if err != nil {
		return PageData{}, err
	}

	userInfo, err := GetUserInfo(ctx, db)
	if err != nil {
		return PageData{}, err
	}

	return PageData{
		Posts:      posts,
		Categories: categories,
		UserInfo:   *userInfo,
	}, nil
}

func GetAllPostRelatedDataByUserID(ctx context.Context, db *repository.Storage, userID int) (PageData, error) {
	posts, err := db.GetAllPostRelatedDataByUser(userID)
	if err != nil {
		return PageData{}, err
	}

	categories, err := GetCategories(db)
	if err != nil {
		return PageData{}, err
	}

	userInfo, err := GetUserInfo(ctx, db)
	if err != nil {
		return PageData{}, err
	}

	return PageData{
		Posts:      posts,
		Categories: categories,
		UserInfo:   *userInfo,
	}, nil
}
