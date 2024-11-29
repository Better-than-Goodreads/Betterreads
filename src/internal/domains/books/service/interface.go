package service

import (
	"errors"

	"github.com/betterreads/internal/domains/books/models"
	er "github.com/betterreads/internal/pkg/errors"
	"github.com/google/uuid"
)

var (
	ErrRatingNotFound      = errors.New("rating not found")
	ErrBookNotFound        = errors.New("book not found")
	ErrPictureNotFound     = errors.New("picture not found")
	ErrRatingAlreadyExists = errors.New("rating already exists")
	ErrReviewAlreadyExists = errors.New("review already exists")
	ErrReviewNotFound      = errors.New("review not found")
	ErrAuthorNotFound      = errors.New("author not found")
	ErrUserNotAuthor       = errors.New("user is not the author")
	ErrUserNotFound        = errors.New("user not found")
	ErrRatingOwnBook       = errors.New("author can't rate his own book")
	ErrDirectionWhenNoSort = errors.New("direction must be empty when sort is empty")

	ErrGenreRequired = er.ErrorParam{
		Name:   "genre",
		Reason: "genre is required",
	}

	ErrGenreNotFound = er.ErrorParam{
		Name:   "genre",
		Reason: "genre not in available genres",
	}

	ErrRatingAmount = er.ErrorParam{
		Name:   "rating",
		Reason: "rating must be between 1 and 5",
	}

	ErrInvalidSort = er.ErrorParam{
		Name:   "sort",
		Reason: "sort must be one of the following: publication_date, total_ratings, avg_ratings",
	}

	ErrInvalidDirection = er.ErrorParam{
		Name:   "direction",
		Reason: "direction must be either 'asc' or 'desc'",
	}
)
var (
	AvailableSorts = []string{"publication_date", "total_ratings", "avg_ratings"}
)

type BooksService interface {
	PublishBook(req *models.NewBookRequest, author uuid.UUID) (*models.BookResponse, error)
	GetBookInfo(bookId uuid.UUID, userId uuid.UUID) (*models.BookResponseWithReview, error)
	GetBooksOfAuthor(authorId uuid.UUID, userId uuid.UUID) ([]*models.BookResponseWithReview, error)
	SearchBooks(name string, genre string, userId uuid.UUID, sort string, isAscDirection string) ([]*models.BookResponseWithReview, error)
	GetBookPicture(id uuid.UUID) ([]byte, error)
	GetBooksInfo(userId uuid.UUID) ([]*models.BookResponseWithReview, error)
	RateBook(bookId uuid.UUID, userId uuid.UUID, rateAmount int) (*models.Rating, error)
	UpdateRating(bookId uuid.UUID, userId uuid.UUID, rateAmount int) error
	GetBookReviews(bookId uuid.UUID) ([]*models.ReviewOfBook, error)
	GetAllReviewsOfUser(userId uuid.UUID) ([]*models.ReviewOfUser, error)
	AddReview(bookId uuid.UUID, userId uuid.UUID, review models.NewReviewRequest) error
	CheckIfUserExists(userId uuid.UUID) bool
	CheckIfAuthorIsRatingOwnBook(bookId uuid.UUID, userId uuid.UUID) (bool, error)
	GetGenres() ([]string, error)
	DeleteReview(bookId uuid.UUID, userId uuid.UUID) error
	DeleteRating(bookId uuid.UUID, userId uuid.UUID) error
	EditReview(bookId uuid.UUID, userId uuid.UUID, review models.NewReviewRequest) error
}
