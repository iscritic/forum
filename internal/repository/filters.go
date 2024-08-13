package repository

import "forum/internal/entity"

func (s *Storage) GetCategoryById(categoryID int) (string, error) {
	var category string

	return category, nil
}

func (s Storage) GetAllPostByCategory(categoryID int) ([]*entity.Post, error) {

	query := `
SELECT p.id, p.title, p.content, p.author_id, p.category_id, p.creation_date
FROM posts p
WHERE p.category_id = ?;
`

	rows, err := s.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CategoryID, &post.CreationDate)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil

}
