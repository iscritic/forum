package service

import "forum/internal/repository"

func FetchPosts(db *repository.Storage) ([]repository.PostRelatedData, error) {

	//TODO: бизнес логика, пагинация страниц

	data, err := db.GetPostsRelatedData()
	if err != nil {
		return nil, err
	}

	//posts, err := db.GetAllPosts()
	//if err != nil {
	//	return posts, err
	//}

	return data, nil
}

func GetPostData(db *repository.Storage, id int) (repository.PostRelatedData, error) {

	post, err := db.GetPostRelatedDataByID(id)
	if err != nil {
		return repository.PostRelatedData{}, err
	}

	return *post, nil
}

func Register(db *repository.Storage, user repository.User) error {

	//TODO бизнес логика комментариев, лайки

	err := db.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}
