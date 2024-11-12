package service

import (
	"errors"
	"github.com/betterreads/internal/domains/bookshelf/models"
	er "github.com/betterreads/internal/pkg/errors"
	"github.com/google/uuid"
)

var (
	ErrBookNotFoundInLibrary = errors.New("book not found")
	ErrBookAlreadyInLibrary  = errors.New("book already in library")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidStatusType     = er.ErrorParam{
		Name:   "status",
		Reason: "status should be: 'plan-to-read', 'reading', 'read' or 'all",
	}
)

type BookshelfService interface {
	GetBookShelf(usedId uuid.UUID, shelfType string) ([]*models.BookInShelfResponse, error)
	AddBookToShelf(userId uuid.UUID, req *models.BookShelfRequest) error
	EditBookInShelf(userId uuid.UUID, req *models.BookShelfRequest) error
	DeleteBookFromShelf(userId uuid.UUID, bookId uuid.UUID) error
}
