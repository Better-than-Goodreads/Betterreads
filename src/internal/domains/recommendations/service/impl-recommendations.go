package service

import (
	"github.com/betterreads/internal/domains/books/models"
	bs "github.com/betterreads/internal/domains/books/service"
	"github.com/betterreads/internal/domains/recommendations/repository"

	"github.com/google/uuid"
)

type RecommendationsServiceImpl struct {
	recommendationsRepository repository.RecommendationsDatabase
	booksService              bs.BooksService
}

func NewRecommendationsServiceImpl(recommendationsRepository repository.RecommendationsDatabase, booksService bs.BooksService) RecommendationsService {
	return &RecommendationsServiceImpl{recommendationsRepository: recommendationsRepository, booksService: booksService}
}

func (rs *RecommendationsServiceImpl) GetRecommendations(userId uuid.UUID) (map[string][]*models.Book, error) {
	userExists := rs.booksService.CheckIfUserExists(userId)
	if !userExists {
		return nil, bs.ErrUserNotFound
	}

	valid := rs.recommendationsRepository.CheckIfUserHasValidShelf(userId)
	if !valid {
		return nil, ErrNeedMoreBooksInShelf
	}

	books, err := rs.recommendationsRepository.GetRecommendations(userId)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (rs *RecommendationsServiceImpl) GetMoreRecommendations(userId uuid.UUID, genre string) ([]*models.Book, error) {
	userExists := rs.booksService.CheckIfUserExists(userId)
	if !userExists {
		return nil, bs.ErrUserNotFound
	}

	valid := rs.recommendationsRepository.CheckIfUserHasValidShelf(userId)
	if !valid {
		return nil, ErrNeedMoreBooksInShelf
	}

	books, err := rs.recommendationsRepository.GetMoreRecommendations(userId, genre)
	if err != nil {
		return nil, err
	}
	return books, nil
}
