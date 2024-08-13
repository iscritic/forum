package repository

import (
	"errors"
	"fmt"
	"forum/internal/entity"
)

func (s *Storage) GetAllPostsRelatedData() ([]entity.PostRelatedData, error) {
	// Получение всех постов
	posts, err := s.GetAllPost()
	if err != nil {
		return nil, err
	}

	// Слайс для хранения всех связанных данных
	var postsRelatedData []entity.PostRelatedData

	// Итерация по каждому посту для получения связанных данных
	for _, post := range posts {
		// Получение связанных комментариев
		comments, err := s.GetAllCommentsRelatedDataByPostID(post.ID)
		if err != nil {
			return nil, err
		}

		// Получение информации об авторе поста
		author, err := s.GetUserByID(post.AuthorID)
		if err != nil {
			return nil, err
		}

		// Получение информации о категории
		categoryName, err := s.GetCategoryById(post.CategoryID)
		if err != nil {
			return nil, err
		}

		post.Likes, post.Dislikes, err = s.GetLikesAndDislikesForPost(post.ID)
		if err != nil {
			return nil, err
		}

		// Заполнение структуры PostRelatedData для текущего поста
		postRelatedData := entity.PostRelatedData{
			Post:     *post,
			CommentR: comments,
			User:     *author,
			Category: entity.Category{ID: post.CategoryID, Name: categoryName},
		}

		// Добавление в слайс
		postsRelatedData = append(postsRelatedData, postRelatedData)
	}

	return postsRelatedData, nil
}

func (s *Storage) GetPostRelatedDataByPostID(postID int) (*entity.PostRelatedData, error) {
	// 1. Fetch the post by ID
	post, err := s.GetPostByID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post: %w", err) // Wrap the error
	}
	if post == nil {
		return nil, errors.New("post not found")
	}

	post.Likes, post.Dislikes, err = s.GetLikesAndDislikesForPost(post.ID)
	if err != nil {
		return nil, err
	}

	// 2. Fetch related comments
	comments, err := s.GetAllCommentsRelatedDataByPostID(postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err) // Wrap the error
	}

	// 3. Fetch author information
	author, err := s.GetUserByID(post.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch author: %w", err) // Wrap the error
	}

	// 4. Fetch category information
	categoryName, err := s.GetCategoryById(post.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch category: %w", err) // Wrap the error
	}

	// 5. Construct and return the PostRelatedData struct
	postRelatedData := &entity.PostRelatedData{
		Post:     *post,
		CommentR: comments,
		User:     *author,
		Category: entity.Category{ID: post.CategoryID, Name: categoryName},
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

	// Получаем уникальные ID авторов
	authorIDs := make(map[int]struct{})
	for _, comment := range comments {
		authorIDs[comment.AuthorID] = struct{}{}
	}

	// Получаем информацию о пользователях
	var authors []entity.User
	for authorID := range authorIDs {
		author, err := s.GetUserByID(authorID)
		if err != nil {
			return nil, err
		}
		authors = append(authors, *author)
	}

	// Создаем отображение ID авторов для быстрого доступа
	authorMap := make(map[int]entity.User)
	for _, author := range authors {
		authorMap[author.ID] = author
	}

	// Создаем результат
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

func (s *Storage) GetAllPostRelatedDataByCategory(categoryID int) ([]entity.PostRelatedData, error) {
	// Получение всех постов
	posts, err := s.GetAllPostByCategory(categoryID)
	if err != nil {
		return nil, err
	}

	fmt.Println("Found some posts:", posts)

	// Слайс для хранения всех связанных данных
	var postsRelatedData []entity.PostRelatedData

	// Итерация по каждому посту для получения связанных данных
	for _, post := range posts {
		// Получение связанных комментариев
		comments, err := s.GetAllCommentsRelatedDataByPostID(post.ID)
		if err != nil {
			return nil, err
		}

		// Получение информации об авторе поста
		author, err := s.GetUserByID(post.AuthorID)
		if err != nil {
			return nil, err
		}

		// Получение информации о категории
		categoryName, err := s.GetCategoryById(post.CategoryID)
		if err != nil {
			return nil, err
		}

		post.Likes, post.Dislikes, err = s.GetLikesAndDislikesForPost(post.ID)
		if err != nil {
			return nil, err
		}

		// Заполнение структуры PostRelatedData для текущего поста
		postRelatedData := entity.PostRelatedData{
			Post:     *post,
			CommentR: comments,
			User:     *author,
			Category: entity.Category{ID: post.CategoryID, Name: categoryName},
		}

		// Добавление в слайс
		postsRelatedData = append(postsRelatedData, postRelatedData)
	}

	fmt.Println("+++++++++++++++", postsRelatedData)

	return postsRelatedData, nil
}
