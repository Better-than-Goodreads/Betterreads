package service

import (
	"errors"
	"github.com/betterreads/internal/domains/books/models"
	"github.com/betterreads/internal/domains/books/repository"
	"github.com/betterreads/internal/domains/books/utils"
)

type BooksService struct {
	booksRepository repository.BooksDatabase
}

func NewBooksService(booksRepository repository.BooksDatabase) *BooksService {
	return &BooksService{booksRepository: booksRepository}
}

func (bs *BooksService) PublishBook(req *models.NewBookRequest) error {
	if len(req.Genres) == 0 {
		return errors.New("at least one genre is required")
	}

	var newBook = utils.MapBookRequestToBookRecord(req)
	if err := bs.booksRepository.SaveBook(newBook); err != nil {
		return err
	}
	return nil
}

func (bs *BooksService) GetBook(name string) (*repository.Book, error) {
	book, err := bs.booksRepository.GetBookByName(name)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (bs *BooksService) RateBook(bookId int, rateAmount int) error {
	err := bs.booksRepository.RateBook(bookId, rateAmount)
	if err != nil {
		return err
	}
	return nil
}