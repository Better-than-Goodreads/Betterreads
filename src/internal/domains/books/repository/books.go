package repository

import (
    "errors"
    "github.com/google/uuid"
	"github.com/betterreads/internal/domains/books/models"
)

var genresDict = map[int]string{
    1:  "Fiction",
    2:  "Non-fiction",
    3:  "Fantasy",
    4:  "Science Fiction",
    5:  "Mystery",
    6:  "Horror",
    7:  "Romance",
    8:  "Thriller",
    9:  "Historical",
    10: "Biography",
    11: "Autobiography",
    12: "Self-help",
    13: "Travel",
    14: "Guide",
    15: "Poetry",
    16: "Drama",
    17: "Satire",
    18: "Anthology",
    19: "Encyclopedia",
    20: "Dictionary",
    21: "Comic",
    22: "Art",
    23: "Cookbook",
}

var (
    ErrGenreNotFound = errors.New("genre not found")
)

type BooksDatabase interface {
	SaveBook(*models.NewBookRequest, string)(*models.Book, error)
	GetBookById(id uuid.UUID) (*models.Book, error)
    GetBooks() ([]*models.Book, error)
	GetBookByName(name string) (*models.Book, error)
	RateBook(bookId uuid.UUID, userId uuid.UUID, rating int) error
	DeleteRating(bookId uuid.UUID, userId uuid.UUID) error
	GetRatings(bookId uuid.UUID, userId uuid.UUID) (*models.Rating, error)
}

