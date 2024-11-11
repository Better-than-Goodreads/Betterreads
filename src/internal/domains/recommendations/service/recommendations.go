package service

import (
	"errors"

	"github.com/betterreads/internal/domains/books/models"
	"github.com/google/uuid"
)

var (
	ErrNeedMoreBooksInShelf = errors.New("Need more than 5 books in shelf")
)

type RecommendationsService interface {
	GetRecommendations(userId uuid.UUID) (map[string][]*models.Book, error)
	GetMoreRecommendations(userId uuid.UUID, genre string) ([]*models.Book, error)
	GetFriendsRecommendations(userId uuid.UUID) ([]*models.Book, error)
}
