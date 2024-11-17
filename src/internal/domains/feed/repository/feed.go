package repository

import (
	"github.com/betterreads/internal/domains/feed/models"
	"github.com/google/uuid"
)

type FeedRepository interface {
	GetFeed(userId uuid.UUID) ([]models.Post, error)
}