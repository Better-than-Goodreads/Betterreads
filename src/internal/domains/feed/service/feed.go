package service

import (
	"github.com/betterreads/internal/domains/feed/models"
	"github.com/google/uuid"
)

type FeedService interface {
	GetFeed(userId uuid.UUID) ([]models.Post, error)
}
