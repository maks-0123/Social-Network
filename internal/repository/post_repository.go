package repository

import (
	"context"
	"my-social-network/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
	pool *pgxpool.Pool
}

func NewPostRepository(pool *pgxpool.Pool) *PostRepository {
	return &PostRepository{pool: pool}
}

// Create создает новый пост
func (r *PostRepository) Create(ctx context.Context, post *models.Post) error {
	query := `
		INSERT INTO posts (content, author_id, parent_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at 
	`

	// Просто передаем post.ParentID - pgx умеет работать с *int
	return r.pool.QueryRow(ctx, query, post.Content, post.AuthorID, post.ParentID).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

// GetFeed возвращает ленту постов
func (r *PostRepository) GetFeed(ctx context.Context, limit int) ([]*models.Post, error) {
	query := `
        SELECT 
            p.id, p.content, p.author_id, p.parent_id, p.created_at, p.updated_at,
            u.id, u.username, u.email, u.avatar_url, u.created_at
        FROM posts p
        JOIN users u ON p.author_id = u.id
        ORDER BY p.created_at DESC
        LIMIT $1
    `

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		var author models.User

		err := rows.Scan(
			&post.ID, &post.Content, &post.AuthorID, &post.ParentID,
			&post.CreatedAt, &post.UpdatedAt,
			&author.ID, &author.Username, &author.Email, &author.AvatarURL, &author.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		post.Author = &author
		posts = append(posts, &post)
	}

	return posts, nil
}
