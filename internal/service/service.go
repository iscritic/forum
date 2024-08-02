package service

import "forum/internal/sqlite"

func FetchPosts(db *sqlite.Storage) ([]*sqlite.Post, error) {

	//TODO: бизнес логика, пагинация страниц

	posts, err := db.GetAllPosts()
	if err != nil {
		return posts, err
	}

	return posts, nil
}

func GetPostData(db *sqlite.Storage, id int) (sqlite.PostData, error) {
	var postData sqlite.PostData

	post, err := db.GetPostByID(id)
	if err != nil {
		return postData, err
	}

	comments, err := db.GetAllComments(id)
	if err != nil {
		return postData, err
	}

	likes, dislikes, err := db.GetLikesAndDislikesForPost(id)
	if err != nil {
		return postData, err
	}

	// Fetch like and dislike counts for each comment
	for i := range comments {
		comments[i].Likes, comments[i].Dislikes, err = db.GetLikesAndDislikesForComment(comments[i].ID)
		if err != nil {
			return postData, err
		}
	}

	postData = sqlite.PostData{
		Post:     *post,
		Comment: comments,
		Likes:    likes,
		Dislikes: dislikes,
	}

	return postData, nil
}

func Register(db *sqlite.Storage, user sqlite.User) error {

	//TODO бизнес логика комментариев, лайки

	err := db.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}
