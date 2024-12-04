package service

import (
	"github.com/betterreads/internal/domains/feed/models"
	"github.com/betterreads/internal/domains/feed/repository"
	us "github.com/betterreads/internal/domains/users/service"
	"github.com/google/uuid"
)

type FeedServiceImpl struct {
	fr repository.FeedRepository
	us us.UsersService
}

func NewFeedServiceImpl(fr repository.FeedRepository, us us.UsersService) FeedService {
	return &FeedServiceImpl{fr: fr, us: us}
}

func (fs *FeedServiceImpl) GetFeed(userId uuid.UUID) ([]models.Post, error) {
	if !fs.us.CheckUserExists(userId) {
		return nil, us.ErrUserNotFound
	}

	posts, err := fs.fr.GetFeed(userId)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
