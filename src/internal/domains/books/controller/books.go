package controller

import (
	"net/http"

	"github.com/betterreads/internal/domains/books/models"
	"github.com/betterreads/internal/domains/books/service"
	"github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"
    "github.com/google/uuid"

)

type BooksController struct {
	bookService *service.BooksService
}

func NewBooksController(bookService *service.BooksService) *BooksController {
	return &BooksController{bookService: bookService}
}

// PublishBook godoc
// @Summary publish a book
// @Description publishes a book 
// @Tags books
// @Accept  json
// @Produce  json
// @Param user body models.NewBookRequest true "Book Request"
// @Success 201 {object} models.Book
// @Failure 400 {object} errors.ErrorDetailsWithParams
// @Failure 500 {object} errors.ErrorDetails
// @Router /books [post]
func (bc *BooksController) PublishBook(ctx *gin.Context) {
    // Validates if the user is an author through jwt
    // isAuthor, _ := ctx.Get("IsAuthor")
    // if isAuthor != true {
    //     errors.SendError(ctx, errors.NewErrNotAuthor())
    // }
    //
	var newBookRequest models.NewBookRequest
	if err := ctx.ShouldBindJSON(&newBookRequest); err != nil {
		errors.SendErrorWithParams(ctx, errors.NewErrParsingRequest(err))
		return
	}
    book , err := bc.bookService.PublishBook(&newBookRequest)
    if err != nil {
		errors.SendError(ctx, errors.NewErrPublishingBook(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"book": book})
}

// GetBook godoc
// @Summary Get book by id 
// @Description Get book id, note that its a UUID
// @Tags books
// @Param id path string true "Book Id"
// @Produce  json
// @Success 200 {object} models.Book
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Router /books/{id} [get]
func (bc *BooksController) GetBook(ctx *gin.Context) {
	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		errors.SendError(ctx, errors.NewErrInvalidBookId(id))
		return
	}

	book, err := bc.bookService.GetBook(uuid)
	if err != nil {
		errors.SendError(ctx, errors.NewErrGettingBook(err))
		return
	}

	if book == nil {
		errors.SendError(ctx, errors.NewErrBookNotFound())
	}

	ctx.JSON(http.StatusCreated, gin.H{"book": book})
}

// GetBooks godoc
// @Summary Get all books
// @Description Get all books
// @Tags books
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.Book
// @Failure 500 {object} errors.ErrorDetails
// @Router /books [get]
func (bc *BooksController) GetBooks(ctx *gin.Context) {
    books, err := bc.bookService.GetBooks()
    if err != nil {
        errors.SendError(ctx, errors.NewErrGettingBooks(err))
        return
    }
    ctx.JSON(http.StatusAccepted, gin.H{"books": books})
}

func (bc *BooksController) RateBook(ctx *gin.Context) {

	var newBookRating models.NewBookRating
	if err := ctx.ShouldBindJSON(&newBookRating); err != nil {
		errors.SendErrorWithParams(ctx, errors.NewErrParsingRequest(err))
		return
	}

	rateAmount := newBookRating.Rating
	bookId := newBookRating.BookId
	userId := newBookRating.UserId

	if err := bc.bookService.RateBook(bookId, userId, rateAmount); err != nil {
		errors.SendError(ctx, errors.NewErrRatingBook(err))
		return
	}

	message := "book rated "

	ctx.JSON(200, gin.H{"message": message,})
}

func (bc *BooksController) DeleteRating(ctx *gin.Context) {

	var newBookRating models.NewBookRating
	if err := ctx.ShouldBindJSON(&newBookRating); err != nil {
		errors.SendErrorWithParams(ctx, errors.NewErrParsingRequest(err))
		return
	}
	bookId := newBookRating.BookId
	userId := newBookRating.UserId

	if err := bc.bookService.DeleteRating(bookId, userId); err != nil {
		//errors.SendError(ctx, errors.NewErrDeletingRating(err))
		return
	}

	ctx.JSON(200, gin.H{"message": "rating deleted"})
}

func (bc *BooksController) GetRatings(ctx *gin.Context) {
	var newBookRating models.NewBookRating
	if err := ctx.ShouldBindJSON(&newBookRating); err != nil {
		errors.SendErrorWithParams(ctx, errors.NewErrParsingRequest(err))
		return
	}
	bookId := newBookRating.BookId
	userId := newBookRating.UserId

	ratings, err := bc.bookService.GetRatings(bookId, userId)
	if err != nil {
		//errors.SendError(ctx, errors.NewErrGettingRatings(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"ratings": ratings})
}
