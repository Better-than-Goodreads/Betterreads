package service

import (
	"errors"
	"fmt"

	"github.com/betterreads/internal/domains/books/models"
	"github.com/betterreads/internal/domains/books/repository"
	"github.com/betterreads/internal/domains/books/utils"
	er "github.com/betterreads/internal/pkg/errors"
	"github.com/google/uuid"
)

var (
	ErrGenreNotFound  = errors.New("genre not found")
	ErrRatingNotFound = errors.New("rating not found")
	ErrBookNotFound   = errors.New("book not found")
    ErrRatingAlreadyExists = errors.New("rating already exists")
    ErrReviewAlreadyExists = errors.New("review already exists")
    ErrReviewNotFound = errors.New("review not found")

	ErrRatingAmount = er.ErrorParam{
		Name:   "rating",
		Reason: "rating must be between 1 and 5",
	}

)

type BooksService struct {
	booksRepository repository.BooksDatabase
}

func NewBooksService(booksRepository repository.BooksDatabase) *BooksService {
	return &BooksService{booksRepository: booksRepository}
}

func (bs *BooksService) PublishBook(req *models.NewBookRequest, author uuid.UUID) (*models.BookResponse, error) {
	if len(req.Genres) == 0 {
		return nil, errors.New("at least one genre is required")
	}

	book, err := bs.booksRepository.SaveBook(req, author)
	if err != nil {
		return nil, err
	}

	bookRes, err := bs.addAuthor(book, book.Author)
	if err != nil {
		return nil, err
	}

	return bookRes, nil
}

func (bs *BooksService) GetBookInfo(bookId uuid.UUID, userId uuid.UUID) (*models.BookResponseWithReview, error) {
	book, err := bs.booksRepository.GetBookById(bookId)
	if err != nil {
		return nil, err
	}

	bookRes, err := bs.mapBookToBookResponseWithReview(book, userId)

    if err != nil {
        return nil, err
    }

	return bookRes, nil
}

func (bs *BooksService) SearchBooksByName(name string, userId uuid.UUID) ([]*models.BookResponseWithReview, error) {
	books, err := bs.booksRepository.GetBooksByName(name)
	if err != nil {
		if errors.Is(err, repository.ErrNoBooksFound) {
			return []*models.BookResponseWithReview{}, nil
		}
		return nil, err
	}

	return bs.mapBooksToBooksResponseWithReview(books, userId)
}

func (bs *BooksService) GetBookPicture(id uuid.UUID) ([]byte, error) {
	book, err := bs.booksRepository.GetBookPictureById(id)
	if err != nil {
		if errors.Is(err, repository.ErrBookNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book picture: %w", err)
	}

	return book, nil
}

func (bs *BooksService) GetBooksInfo(userId uuid.UUID) ([]*models.BookResponseWithReview, error) {
	books, err := bs.booksRepository.GetBooks()
	if err != nil {
		return nil, err
	}

	return bs.mapBooksToBooksResponseWithReview(books, userId)
}

func (bs *BooksService) mapBooksToBooksResponseWithReview(books []*models.Book, userId uuid.UUID) ([]*models.BookResponseWithReview, error) {
	booksResponses := []*models.BookResponseWithReview{}

	for _, book := range books {
		bookResponse , err:= bs.mapBookToBookResponseWithReview(book, userId)
        if err != nil {
            return nil, err
        }
		booksResponses = append(booksResponses, bookResponse)
	}
	return booksResponses, nil
}

func (bs *BooksService) mapBookToBookResponseWithReview(book *models.Book, userId uuid.UUID) (*models.BookResponseWithReview, error) {
    var err error
    bookRes := &models.BookResponseWithReview{}
    if userId != uuid.Nil {
        bookRes.Review, err = bs.booksRepository.GetBookReviewOfUser(book.Id, userId)
        if err != nil {
            if errors.Is(err, repository.ErrRatingNotFound) {
                bookRes.Review = nil
            } else {
                return nil, err
            }
        }
    }

    bookRes.Book, err = bs.addAuthor(book, book.Author)
    if err != nil {
        return nil, err
    }

    return bookRes, nil
}

func (bs *BooksService) RateBook(bookId uuid.UUID, userId uuid.UUID, rateAmount int) (*models.Rating, error) {
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
	
	bookRating, err := bs.booksRepository.RateBook(bookId, userId, rateAmount)
	if err != nil {
		return nil, err
	}
	return bookRating, nil
}

func (bs *BooksService) UpdateRating(bookId uuid.UUID, userId uuid.UUID, rateAmount int) (error){
    if rateAmount < 1 || rateAmount > 5 {
        return  ErrRatingAmount
    }

    if exists, err := bs.booksRepository.CheckIfRatingExists(bookId, userId); err != nil {
		return  err
	} else if !exists {
        return  ErrRatingNotFound
	}

    err := bs.booksRepository.UpdateRating(bookId, userId, rateAmount)
    if err != nil {
        return err
    }

    return nil
}

// func (bs *BooksService) DeleteRating(bookId uuid.UUID, userId uuid.UUID) error {
// 	err := bs.booksRepository.DeleteRating(bookId, userId)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
//
func (bs *BooksService) GetBookReviews(bookId uuid.UUID) ([]*models.Review, error){
    reviews, err := bs.booksRepository.GetBookReviews(bookId)
    if err != nil {
        return nil, err
    }
    return reviews, nil
}

func (bs *BooksService) addAuthor(book *models.Book, author uuid.UUID) (*models.BookResponse, error) {
	author_name, err := bs.booksRepository.GetAuthorName(author)
	if err != nil {
		return nil, err
	}
	bookRes := utils.MapBookToBookResponse(book, author_name)
	return bookRes, nil
}

func (bs *BooksService) AddReview(bookId uuid.UUID, userId uuid.UUID, review models.NewReviewRequest) error {
	if review.Rating < 1 || review.Rating > 5 {
		return ErrRatingAmount
	}

    bookExists := bs.booksRepository.CheckIfBookExists(bookId) 
    if !bookExists {
        return ErrBookNotFound
    }
    exists , err := bs.booksRepository.CheckifReviewExists(bookId, userId)
    if err != repository.ErrReviewEmpty && err != nil {
        return err
    }
    if exists {
        return ErrReviewAlreadyExists
    }
    
    
    if err == repository.ErrReviewEmpty {
        err = bs.booksRepository.EditReview(bookId, userId, review.Rating , review.Review)
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
