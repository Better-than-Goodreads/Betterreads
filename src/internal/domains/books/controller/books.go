package controller

import (
	"net/http"

	"github.com/betterreads/internal/domains/books/models"
	"github.com/betterreads/internal/domains/books/service"
	"github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"



)

type BooksController struct {
	bookService *service.BooksService
}

func NewBooksController(bookService *service.BooksService) *BooksController {
	return &BooksController{bookService: bookService}
}

func (bc *BooksController) PublishBook(ctx *gin.Context) {
	var newBookRequest models.NewBookRequest
	if err := ctx.ShouldBindJSON(&newBookRequest); err != nil {
		errors.SendErrorWithParams(ctx, errors.NewErrParsingRequest(err))
		return
	}

	if err := bc.bookService.PublishBook(&newBookRequest); err != nil {
		errors.SendError(ctx, errors.NewErrPublishingBook(err))
		return
	}

	ctx.JSON(200, gin.H{"message": "book published"})
}

func (bc *BooksController) GetBook(ctx *gin.Context) {

	bookName := ctx.Param("book-name")
	book, err := bc.bookService.GetBook(bookName)
	if err != nil {
		errors.SendError(ctx, errors.NewErrGettingBook(err))
		return
	}

	if book == nil {
		errors.SendError(ctx, errors.NewErrBookNotFound())
	}

	ctx.JSON(http.StatusCreated, gin.H{"book": book})
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