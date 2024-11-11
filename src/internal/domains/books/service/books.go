package service

import (
	"errors"
	"fmt"

	"github.com/betterreads/internal/domains/books/models"
	"github.com/betterreads/internal/domains/books/repository"
	"github.com/betterreads/internal/domains/books/utils"
	"github.com/google/uuid"
)

type BooksServiceImpl struct {
	booksRepository repository.BooksDatabase
}

func NewBooksServiceImpl(booksRepository repository.BooksDatabase) BooksService {
	return &BooksServiceImpl{booksRepository: booksRepository}
}

func (bs *BooksServiceImpl) PublishBook(req *models.NewBookRequest, author uuid.UUID) (*models.BookResponse, error) {
	if len(req.Genres) == 0 {
		return nil, ErrGenreRequired
	}

	if !bs.booksRepository.CheckIfUserExists(author) {
		return nil, ErrAuthorNotFound
	}

	if !bs.booksRepository.CheckIfUserIsAuthor(author) {
		return nil, ErrUserNotAuthor
	}

	book, err := bs.booksRepository.SaveBook(req, author)
	if err != nil {
		if errors.Is(err, repository.ErrGenreNotFound) {
			return nil, ErrGenreNotFound
		}
		return nil, err
	}

	res := utils.MapBookToBookResponse(book)

	return res, nil
}

func (bs *BooksServiceImpl) GetBookInfo(bookId uuid.UUID, userId uuid.UUID) (*models.BookResponseWithReview, error) {
	book, err := bs.booksRepository.GetBookById(bookId)
	if err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	bookRes, err := bs.mapBookToBookResponseWithReview(book, userId)

	if err != nil {
		return nil, err
	}

	return bookRes, nil
}

func (bs *BooksServiceImpl) GetBooksOfAuthor(authorId uuid.UUID, userId uuid.UUID) ([]*models.BookResponseWithReview, error) {
	exists := bs.booksRepository.CheckIfUserExists(authorId)
	if !exists {
		return nil, ErrAuthorNotFound
	}

	isAuthor := bs.booksRepository.CheckIfUserIsAuthor(authorId)
	if !isAuthor {
		return nil, ErrUserNotAuthor
	}

	books, err := bs.booksRepository.GetBooksOfAuthor(authorId)
	if err != nil {
		if errors.Is(err, repository.ErrAuthorNotFound) {
			return nil, ErrAuthorNotFound
		}
		return nil, err
	}

	return bs.mapBooksToBooksResponseWithReview(books, userId)
}

func (bs *BooksServiceImpl) SearchBooks(name string, genre string, userId uuid.UUID, sort string, direction string) ([]*models.BookResponseWithReview, error) {
	if sort != "" {
		if err := validateSort(sort); err != nil {
			return nil, err
		}
	}

	if sort == "" && direction != "" {
		return nil, ErrDirectionWhenNoSort
	}

	if direction != "" && direction != "asc" && direction != "desc" {
		return nil, ErrInvalidDirection
	}

	isDirAsc := direction == "asc"

	books, err := bs.booksRepository.GetBooksByNameAndGenre(name, genre, sort, isDirAsc)
	if err != nil {
		if errors.Is(err, repository.ErrGenreNotFound) {
			return nil, ErrGenreNotFound
		}
		return nil, err
	}

	return bs.mapBooksToBooksResponseWithReview(books, userId)
}

func (bs *BooksServiceImpl) GetBookPicture(id uuid.UUID) ([]byte, error) {
	exists := bs.booksRepository.CheckIfBookExists(id)
	if !exists {
		return nil, ErrBookNotFound
	}

	book, err := bs.booksRepository.GetBookPictureById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get book picture: %w", err)
	}

	return book, nil
}

func (bs *BooksServiceImpl) GetBooksInfo(userId uuid.UUID) ([]*models.BookResponseWithReview, error) {
	books, err := bs.booksRepository.GetBooks()
	if err != nil {
		return nil, err
	}
	res, err := bs.mapBooksToBooksResponseWithReview(books, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (bs *BooksServiceImpl) mapBooksToBooksResponseWithReview(books []*models.Book, userId uuid.UUID) ([]*models.BookResponseWithReview, error) {
	booksResponses := []*models.BookResponseWithReview{}

	for _, book := range books {
		bookResponse, err := bs.mapBookToBookResponseWithReview(book, userId)
		if err != nil {
			return nil, err
		}
		booksResponses = append(booksResponses, bookResponse)
	}
	return booksResponses, nil
}

func (bs *BooksServiceImpl) mapBookToBookResponseWithReview(book *models.Book, userId uuid.UUID) (*models.BookResponseWithReview, error) {
	var err error
	bookRes := &models.BookResponseWithReview{}
	if userId != uuid.Nil {
		bookRes.Review, err = bs.booksRepository.GetBookReviewOfUser(book.Id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrReviewNotFound) {
				bookRes.Review = nil
			} else {
				return nil, err
			}
		}
		bookRes.BookShelfStatus, err = bs.booksRepository.GetBookshelfStatusOfUser(book.Id, userId)
		if err != nil {
			if errors.Is(err, repository.ErrBookNotInShelf) {
				bookRes.BookShelfStatus = nil
			} else {
				return nil, err
			}
		}
	}

	bookRes.Book = utils.MapBookToBookResponse(book)

	return bookRes, nil
}

func (bs *BooksServiceImpl) RateBook(bookId uuid.UUID, userId uuid.UUID, rateAmount int) (*models.Rating, error) {
	if rateAmount < 1 || rateAmount > 5 {
		return nil, ErrRatingAmount
	}

	bookExists := bs.booksRepository.CheckIfBookExists(bookId)
	if !bookExists {
		return nil, ErrBookNotFound
	}

	if exists, err := bs.booksRepository.CheckIfRatingExists(bookId, userId); err != nil {
		return nil, err
	} else if exists {
		return nil, ErrRatingAlreadyExists
	}

	if ratingOwnBook, err := bs.CheckIfAuthorIsRatingOwnBook(bookId, userId); err != nil {
		return nil, err
	} else if ratingOwnBook {
		return nil, ErrRatingOwnBook
	}

	bookRating, err := bs.booksRepository.RateBook(bookId, userId, rateAmount)
	if err != nil {
		return nil, err
	}
	return bookRating, nil
}

func (bs *BooksServiceImpl) CheckIfAuthorIsRatingOwnBook(bookId uuid.UUID, userId uuid.UUID) (bool, error) {
	isAuthor := bs.booksRepository.CheckIfUserIsAuthor(userId)
	if isAuthor {
		AuthorsBooks, err := bs.booksRepository.GetBooksOfAuthor(userId)
		if err != nil {
			return false, err
		}
		for _, book := range AuthorsBooks {
			if book.Id == bookId {
				return true, nil
			}
		}
	}
	return false, nil
}

func (bs *BooksServiceImpl) UpdateRating(bookId uuid.UUID, userId uuid.UUID, rateAmount int) error {
	if rateAmount < 1 || rateAmount > 5 {
		return ErrRatingAmount
	}

	if exists, err := bs.booksRepository.CheckIfRatingExists(bookId, userId); err != nil {
		return err
	} else if !exists {
		return ErrRatingNotFound
	}

	if ratingOwnBook, err := bs.CheckIfAuthorIsRatingOwnBook(bookId, userId); err != nil {
		return err
	} else if ratingOwnBook {
		return ErrRatingOwnBook
	}

	err := bs.booksRepository.UpdateRating(bookId, userId, rateAmount)
	if err != nil {
		return err
	}
	return nil
}

func (bs *BooksServiceImpl) GetBookReviews(bookId uuid.UUID) ([]*models.ReviewOfBook, error) {
	reviews, err := bs.booksRepository.GetBookReviews(bookId)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (bs *BooksServiceImpl) GetAllReviewsOfUser(userId uuid.UUID) ([]*models.ReviewOfUser, error) {
	exists := bs.booksRepository.CheckIfUserExists(userId)
	if !exists {
		return nil, ErrUserNotFound
	}

	reviews, err := bs.booksRepository.GetAllReviewsOfUser(userId)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (bs *BooksServiceImpl) AddReview(bookId uuid.UUID, userId uuid.UUID, review models.NewReviewRequest) error {
	if review.Rating < 1 || review.Rating > 5 {
		return ErrRatingAmount
	}

	bookExists := bs.booksRepository.CheckIfBookExists(bookId)
	if !bookExists {
		return ErrBookNotFound
	}

	exists, err := bs.booksRepository.CheckifReviewExists(bookId, userId)
	if err != repository.ErrReviewEmpty && err != nil {
		return err
	}

	if exists {
		return ErrReviewAlreadyExists
	}

	if ratingOwnBook, err := bs.CheckIfAuthorIsRatingOwnBook(bookId, userId); err != nil {
		return err
	} else if ratingOwnBook {
		return ErrRatingOwnBook
	}

	if err == repository.ErrReviewEmpty {
		err = bs.booksRepository.EditReview(bookId, userId, review.Rating, review.Review)
		if err != nil {
			return err
		}
	} else {
		err = bs.booksRepository.AddReview(bookId, userId, review.Review, review.Rating)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bs *BooksServiceImpl) CheckIfUserExists(userId uuid.UUID) bool {
	return bs.booksRepository.CheckIfUserExists(userId)
}

func validateSort(sort string) error {
	for _, s := range AvailableSorts {
		if s == sort {
			return nil
		}
	}
	return ErrInvalidSort
}

func (bs *BooksServiceImpl) GetGenres() ([]string, error) {
	genres, err := bs.booksRepository.GetGenres()
	if err != nil {
		return nil, err
	}
	return genres, nil
}
