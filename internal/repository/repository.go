package repository

import (
	"errors"
	"fmt"

	"forum/internal/entity"
)

func (s *Storage) GetAllPostsRelatedData() ([]*entity.PostRelatedData, error) {
	posts, err := s.GetAllPost()
	if err != nil {
		return nil, err
	}

	postsRelatedData, err := s.postAdoptionCenter(posts)
	if err != nil {
		return nil, err
	}

	return postsRelatedData, nil
}

func (s *Storage) GetPostRelatedDataByPostID(postID int) (*entity.PostRelatedData, error) {
	post, err := s.GetPostByID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post: %w", err)
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	post.Likes, post.Dislikes, err = s.GetLikesAndDislikesForPost(post.ID)
	if err != nil {
		return nil, err
	}

	comments, err := s.GetAllCommentsRelatedDataByPostID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}

	author, err := s.GetUserByID(post.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch author: %w", err)
	}

	category, err := s.GetCategoryById(post.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category: %w", err)
	}

	postRelatedData := &entity.PostRelatedData{
		Post:     *post,
		CommentR: comments,
		User:     *author,
		Category: *category,
	}

	return postRelatedData, nil
}

func (s *Storage) GetAllCommentsRelatedDataByPostID(postID int) ([]entity.CommentRelatedData, error) {
	rows, err := s.db.Query(`
        SELECT id, post_id, content, author_id, likes, dislikes, creation_date 
        FROM comments 
        WHERE post_id = ? 
        ORDER BY creation_date DESC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []entity.Comment
	for rows.Next() {
		var comment entity.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Content,
			&comment.AuthorID,
			&comment.Likes,
			&comment.Dislikes,
			&comment.CreationDate,
		); err != nil {
			return nil, err
		}

		comment.Likes, comment.Dislikes, err = s.GetLikesAndDislikesForComment(comment.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	authorIDs := make(map[int]struct{})
	for _, comment := range comments {
		authorIDs[comment.AuthorID] = struct{}{}
	}

	var authors []entity.User
	for authorID := range authorIDs {
		author, err := s.GetUserByID(authorID)
		if err != nil {
			return nil, err
		}
		authors = append(authors, *author)
	}
	authorMap := make(map[int]entity.User)
	for _, author := range authors {
		authorMap[author.ID] = author
	}

	var commentsRelatedData []entity.CommentRelatedData
	for _, comment := range comments {
		user := authorMap[comment.AuthorID]
		commentsRelatedData = append(commentsRelatedData, entity.CommentRelatedData{
			Comment: comment,
			User:    user,
		})
	}

	return commentsRelatedData, nil
}

func (s *Storage) GetAllPostRelatedDataByCategory(categoryID int) ([]*entity.PostRelatedData, error) {
	posts, err := s.GetAllPostByCategory(categoryID)
	if err != nil {
		return nil, err
	}

	postsRelatedData, err := s.postAdoptionCenter(posts)
	if err != nil {
		return nil, err
	}

	return postsRelatedData, nil
}

func (s *Storage) GetAllPostRelatedDataByUser(userID int) ([]*entity.PostRelatedData, error) {
	posts, err := s.GetAllPostByUser(userID)
	if err != nil {
		return nil, err
	}

	postsRelatedData, err := s.postAdoptionCenter(posts)
	if err != nil {
		return nil, err
	}
	return postsRelatedData, nil
}

func (s *Storage) GetMyLikedPosts(userID int) ([]*entity.PostRelatedData, error) {
	posts, err := s.GetAllLikedPosts(userID)
	if err != nil {
		return nil, err
	}

	postsRelatedData, err := s.postAdoptionCenter(posts)
	if err != nil {
		return nil, err
	}

	return postsRelatedData, nil
}

func (s *Storage) postAdoptionCenter(posts []*entity.Post) ([]*entity.PostRelatedData, error) {
	var postsRelatedData []*entity.PostRelatedData

	for _, post := range posts {
		comments, err := s.GetAllCommentsRelatedDataByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		author, err := s.GetUserByID(post.AuthorID)
		if err != nil {
			return nil, err
		}
		category, err := s.GetCategoryById(post.CategoryID)
		if err != nil {
			return nil, err
		}
		post.Likes, post.Dislikes, err = s.GetLikesAndDislikesForPost(post.ID)
		if err != nil {
			return nil, err
		}

		postRelatedData := &entity.PostRelatedData{
			Post:     *post,
			CommentR: comments,
			User:     *author,
			Category: *category,
		}
		postsRelatedData = append(postsRelatedData, postRelatedData)
	}

	return postsRelatedData, nil
}
