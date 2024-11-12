package service

import (
	_ "errors"
	bs "github.com/betterreads/internal/domains/books/service"
	"github.com/betterreads/internal/domains/bookshelf/models"
	"github.com/betterreads/internal/domains/bookshelf/repository"
	"github.com/google/uuid"
)

type BookShelfServiceImpl struct {
	r           repository.BookshelfDatabase
	bookService bs.BooksService
}

func NewBookShelfServiceImpl(r repository.BookshelfDatabase, bs bs.BooksService) BookshelfService {
	return &BookShelfServiceImpl{r: r, bookService: bs}
}

func (bs *BookShelfServiceImpl) GetBookShelf(userId uuid.UUID, shelfType string) ([]*models.BookInShelfResponse, error) {
	userExists := bs.bookService.CheckIfUserExists(userId)
	if !userExists {
		return nil, ErrUserNotFound
	}

	status := models.BookShelfType(shelfType)
	if !validate_status(status) {
		return nil, ErrInvalidStatusType
	}

	bookShelf, err := bs.r.GetBookShelf(userId, status)
	if err != nil {
		return nil, err
	}

	return bookShelf, nil
}

func (bs *BookShelfServiceImpl) AddBookToShelf(userId uuid.UUID, req *models.BookShelfRequest) error {
	userExists := bs.bookService.CheckIfUserExists(userId)
	if !userExists {
		return ErrUserNotFound
	}

	status := models.BookShelfType(req.Status)
	if !validate_status(status) {
		return ErrInvalidStatusType
	}

	exists := bs.r.CheckIfBookIsInUserShelf(userId, req.BookId)
	if exists {
		return ErrBookAlreadyInLibrary
	}

	err := bs.r.AddBookToShelf(userId, req)
	if err != nil {
		return err
	}

	return nil
}

func (bs *BookShelfServiceImpl) EditBookInShelf(userId uuid.UUID, req *models.BookShelfRequest) error {
	userExists := bs.bookService.CheckIfUserExists(userId)
	if !userExists {
		return ErrUserNotFound
	}

	exits := bs.r.CheckIfBookIsInUserShelf(userId, req.BookId)
	if !exits {
		return ErrBookNotFoundInLibrary
	}

	status := models.BookShelfType(req.Status)
	if !validate_status(status) {
		return ErrInvalidStatusType
	}

	err := bs.r.EditBookInShelf(userId, req)
	if err != nil {
		return err
	}

	return nil
}

func validate_status(status models.BookShelfType) bool {
	for _, s := range models.ValidBookShelfTypes {
		if s == status {
			return true
		}
	}
	return false
}


func (bs *BookShelfServiceImpl) DeleteBookFromShelf(userId uuid.UUID, bookId uuid.UUID) error {
	userExists := bs.bookService.CheckIfUserExists(userId)
	if !userExists {
		return ErrUserNotFound
	}

	exits := bs.r.CheckIfBookIsInUserShelf(userId, bookId)
	if !exits {
		return ErrBookNotFoundInLibrary
	}

	err := bs.r.DeleteBookFromShelf(userId, bookId)
	if err != nil {
		return err
	}

	return nil
}