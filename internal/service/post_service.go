package service

import (
	"context"
	"my-social-network/internal/models"
	"my-social-network/internal/repository"
)

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
	return &PostService{postRepo: postRepo}
}
func (s *PostService) CreatePost(userID int, content string) error {
	ctx := context.Background()

	post := &models.Post{
		Content:  content,
		AuthorID: userID,
		ParentID: nil,
	}
	return s.postRepo.Create(ctx, post)
}
func (s *PostService) GetFeed(limit int) ([]*models.Post, error) {
	ctx := context.Background()
	return s.postRepo.GetFeed(ctx, limit)
}
