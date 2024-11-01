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
    ErrRatingNotFound = errors.New("review not found")
    ErrAuthorNotFound = errors.New("author not found")
    ErrBookNotFound = errors.New("book not found")
    ErrNoBooksFound = errors.New("no books found")
    ErrRatingAlreadyExists = errors.New("rating already exists")
    ErrReviewAlreadyExists = errors.New("review already exists")
    ErrReviewNotFound = errors.New("review not found")
    ErrReviewEmpty = errors.New("review is empty")
    ErrUserNotFound = errors.New("user not found")
)

type BooksDatabase interface {
	SaveBook(*models.NewBookRequest, uuid.UUID)(*models.Book, error)
	GetBookById(id uuid.UUID) (*models.Book, error)
    GetBookPictureById(id uuid.UUID) ([]byte, error)
    GetBooks() ([]*models.Book, error)

    // GetBooksOfAuthor returns all books of an author, if it doesn't exist returns ErrAuthorNotFound
    GetBooksOfAuthor(authorId uuid.UUID) ([]*models.Book, error)
    GetBooksByName(name string) ([]*models.Book, error)
    CheckIfBookExists(bookId uuid.UUID) bool
    // RATE
	RateBook(bookId uuid.UUID, userId uuid.UUID, rating int) (*models.Rating, error)
    UpdateRating(bookId uuid.UUID, userId uuid.UUID, rating int) (error)
    CheckIfRatingExists(bookId uuid.UUID, userId uuid.UUID) (bool, error)
    GetBookReviews(bookID uuid.UUID) ([]*models.Review, error)
	// DeleteRating(bookId uuid.UUID, userId uuid.UUID) error

    CheckIfAuthorExists(authorId uuid.UUID) bool

    AddReview(bookId uuid.UUID, userId uuid.UUID, review string, rating int) error
    CheckifReviewExists(bookId uuid.UUID, userId uuid.UUID) (bool, error)
	GetBookReviewOfUser(bookId uuid.UUID, userId uuid.UUID) (*models.Review, error)

    // GetAllReviewsOfUser returns all reviews of a user, if it doesn't exist returns ErrUserNotFound
    GetAllReviewsOfUser(userId uuid.UUID) ([]*models.Review, error)
    EditReview(bookId uuid.UUID, userId uuid.UUID, rating int, review string) (error)
}

