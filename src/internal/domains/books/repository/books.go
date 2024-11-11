package repository

import (
	"errors"
	"github.com/betterreads/internal/domains/books/models"
	"github.com/google/uuid"
)

var GenresDict = map[int]string{
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
	ErrGenreNotFound       = errors.New("genre not found")
	ErrRatingNotFound      = errors.New("review not found")
	ErrAuthorNotFound      = errors.New("author not found")
	ErrBookNotFound        = errors.New("book not found")
	ErrNoBooksFound        = errors.New("no books found")
	ErrRatingAlreadyExists = errors.New("rating already exists")
	ErrReviewAlreadyExists = errors.New("review already exists")
	ErrBookNotInShelf      = errors.New("book not in shelf")
	ErrReviewNotFound      = errors.New("review not found")
	ErrReviewEmpty         = errors.New("review is empty")
	ErrUserNotFound        = errors.New("user not found")
)

type BooksDatabase interface {
	SaveBook(*models.NewBookRequest, uuid.UUID) (*models.Book, error)
	GetBookById(id uuid.UUID) (*models.Book, error)
	GetBookPictureById(id uuid.UUID) ([]byte, error)
	GetBooks() ([]*models.Book, error)
	GetBooksOfAuthor(authorId uuid.UUID) ([]*models.Book, error)
	GetBooksByNameAndGenre(name string, genre string, sort string, directAsc bool) ([]*models.Book, error)
	GetGenresForBook(book_id uuid.UUID) ([]string, error)
    GetGenres() ([]string, error)

	CheckIfBookExists(bookId uuid.UUID) bool
	CheckIfUserExists(userId uuid.UUID) bool
	CheckIfUserIsAuthor(authorId uuid.UUID) bool

	RateBook(bookId uuid.UUID, userId uuid.UUID, rating int) (*models.Rating, error)
	UpdateRating(bookId uuid.UUID, userId uuid.UUID, rating int) error
	CheckIfRatingExists(bookId uuid.UUID, userId uuid.UUID) (bool, error)

	AddReview(bookId uuid.UUID, userId uuid.UUID, review string, rating int) error
	CheckifReviewExists(bookId uuid.UUID, userId uuid.UUID) (bool, error)
	GetBookReviews(bookID uuid.UUID) ([]*models.ReviewOfBook, error)
	GetBookReviewOfUser(bookId uuid.UUID, userId uuid.UUID) (*models.Review, error)
	GetBookshelfStatusOfUser(bookId uuid.UUID, userId uuid.UUID) (*string, error)
	GetAllReviewsOfUser(userId uuid.UUID) ([]*models.ReviewOfUser, error)
	EditReview(bookId uuid.UUID, userId uuid.UUID, rating int, review string) error
}
