package service

import (
	"errors"

	"github.com/betterreads/internal/domains/books/models"
	"github.com/betterreads/internal/domains/books/repository"
	"github.com/betterreads/internal/domains/books/utils"
	"github.com/google/uuid"
)

type BooksService struct {
	booksRepository repository.BooksDatabase
}

func NewBooksService(booksRepository repository.BooksDatabase) *BooksService {
	return &BooksService{booksRepository: booksRepository}
}

func (bs *BooksService) PublishBook(req *models.NewBookRequest, author string) (*models.BookResponse, error) {
	if len(req.Genres) == 0 {
		return nil, errors.New("at least one genre is required")
    }

    book ,err := bs.booksRepository.SaveBook(req, author)
    if err != nil {
       return nil, err
    }

    bookRes := utils.MapBookToBookResponse(book)
	return bookRes, nil
}

func (bs *BooksService) GetBook(id uuid.UUID) (*models.Book, error) {
	book, err := bs.booksRepository.GetBookById(id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (bs *BooksService) GetBooks() ([]*models.Book, error) {
    books, err := bs.booksRepository.GetBooks()
    if err != nil {
        return nil, err
    }
    return books, nil
}

func (bs *BooksService) RateBook(bookId uuid.UUID, userId uuid.UUID, rateAmount int) error {

	if rateAmount < 1 || rateAmount > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	err := bs.booksRepository.RateBook(bookId, userId, rateAmount)
	if err != nil {
		return err
	}
	return nil
}

func (bs *BooksService) DeleteRating(bookId uuid.UUID, userId uuid.UUID) error {
	
	err := bs.booksRepository.DeleteRating(bookId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (bs *BooksService) GetRatings(bookId uuid.UUID, userId uuid.UUID) (int, error) {
	rating, err := bs.booksRepository.GetRatings(bookId, userId)
	if err != nil {
		return -1, err
	}


	return rating.Rating, nil
}
