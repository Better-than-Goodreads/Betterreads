package repository

import (
	"errors"
	"github.com/betterreads/internal/domains/books/models"
	"github.com/google/uuid"
)

var (
	ErrNeedMoreBooksInShelf = errors.New("need 5 or more books in shelf")
)

type RecommendationsDatabase interface {
	GetMoreRecommendations(userId uuid.UUID, genre string) ([]*models.Book, error)
	GetRecommendations(userId uuid.UUID) (map[string][]*models.Book, error)
	CheckIfUserHasValidShelf(userId uuid.UUID) bool
	GetFriendsRecommendations(userId uuid.UUID) ([]*models.Book, error)
}
