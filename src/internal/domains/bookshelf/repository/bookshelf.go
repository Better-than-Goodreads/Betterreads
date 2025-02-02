package repository

import (
	"errors"

	"github.com/betterreads/internal/domains/bookshelf/models"
	er "github.com/betterreads/internal/pkg/errors"
	"github.com/google/uuid"
)

var (
	ErrBookNotFoundInLibrary = errors.New("book not found")
	ErrBookAlreadyInLibrary  = errors.New("book already in library")
	ErrBookNotInLibrary      = errors.New("book not in library")
	ErrGenreNotFound         = errors.New("genre not found")

	ErrInvaliStatusType = er.ErrorParam{
		Name:   "status",
		Reason: "status should be: 'Plan-To-Read', 'Reading' or 'Read'",
	}
)

type BookshelfDatabase interface {
	GetBookShelf(usedId uuid.UUID, ShelfType models.BookShelfType) ([]*models.BookInShelfResponse, error)
	AddBookToShelf(userId uuid.UUID, req *models.BookShelfRequest) error
	EditBookInShelf(userId uuid.UUID, req *models.BookShelfRequest) error
	SearchBookShelf(userId uuid.UUID, shelfType models.BookShelfType, genre string, sort string, isDirAsc bool) ([]*models.BookInShelfResponse, error)
	CheckIfBookIsInUserShelf(userId uuid.UUID, bookId uuid.UUID) bool
	DeleteBookFromShelf(userId uuid.UUID, bookId uuid.UUID) error
}
