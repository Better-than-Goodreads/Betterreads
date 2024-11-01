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

func NewBooksServiceImpl(booksRepository repository.BooksDatabase) BooksService{
	return &BooksServiceImpl{booksRepository: booksRepository}
}

func (bs *BooksServiceImpl) PublishBook(req *models.NewBookRequest, author uuid.UUID) (*models.BookResponse, error) {
	if len(req.Genres) == 0 {
		return nil, ErrGenreRequired
    }

    if !bs.booksRepository.CheckIfAuthorExists(author) {
        return nil, ErrAuthorNotFound
    }

	book, err := bs.booksRepository.SaveBook(req, author)
	if err != nil {
        if errors.Is(err, repository.ErrGenreNotFound){
            return nil, ErrGenreNotFound
        }
		return nil, err
	}

	bookRes, err := bs.addAuthor(book, book.Author)
	if err != nil {
		return nil, err
	}

	return bookRes, nil
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
	books, err := bs.booksRepository.GetBooksOfAuthor(authorId)
	if err != nil {
		if errors.Is(err, repository.ErrAuthorNotFound) {
			return nil, ErrAuthorNotFound
		}
		return nil, err
	}

	return bs.mapBooksToBooksResponseWithReview(books, userId)
}

func (bs *BooksServiceImpl) SearchBooksByName(name string, userId uuid.UUID) ([]*models.BookResponseWithReview, error) {
	books, err := bs.booksRepository.GetBooksByName(name)
	if err != nil {
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

	return bs.mapBooksToBooksResponseWithReview(books, userId)
}

func (bs *BooksServiceImpl) mapBooksToBooksResponseWithReview(books []*models.Book, userId uuid.UUID) ([]*models.BookResponseWithReview, error) {
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
    }

    bookRes.Book, err = bs.addAuthor(book, book.Author)
    if err != nil {
        return nil, err
    }

    return bookRes, nil
}

func (bs *BooksServiceImpl) RateBook(bookId uuid.UUID, userId uuid.UUID, rateAmount int) (*models.Rating, error){
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

func (bs *BooksServiceImpl) UpdateRating(bookId uuid.UUID, userId uuid.UUID, rateAmount int) (error){
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

func (bs *BooksServiceImpl) GetBookReviews(bookId uuid.UUID) ([]*models.Review, error){
    reviews, err := bs.booksRepository.GetBookReviews(bookId)
    if err != nil {
        return nil, err
    }
    return reviews, nil
}

func (bs *BooksServiceImpl) GetAllReviewsOfUser(userId uuid.UUID) ([]*models.Review, error){
	reviews, err := bs.booksRepository.GetAllReviewsOfUser(userId)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return reviews, nil
}

func (bs *BooksServiceImpl) addAuthor(book *models.Book, author uuid.UUID) (*models.BookResponse, error) {
	author_name, err := bs.booksRepository.GetAuthorName(author)
	if err != nil {
        if errors.Is(err, repository.ErrAuthorNotFound) {
            return nil, ErrAuthorNotFound
        }
		return nil, err
	}
	bookRes := utils.MapBookToBookResponse(book, author_name)
	return bookRes, nil
}

func (bs *BooksServiceImpl) AddReview(bookId uuid.UUID, userId uuid.UUID, review models.NewReviewRequest) error {
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
