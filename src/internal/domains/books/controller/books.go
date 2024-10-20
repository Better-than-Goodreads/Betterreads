package controller

import (
	"net/http"

	"github.com/betterreads/internal/domains/books/models"
	"github.com/betterreads/internal/domains/books/service"
	"github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"

	"strconv"

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

	strBookId := ctx.Param("book-id")
	strRateAmount := ctx.Param("rate-amount")

	bookId, err := strconv.Atoi(strBookId)
	if err != nil {
		errors.SendError(ctx, errors.NewErrInvalidBookId(strBookId))
		return
	}
	rateAmount, err := strconv.Atoi(strRateAmount)
	if err != nil {
		errors.SendError(ctx, errors.NewErrInvalidRating(strRateAmount))
		return
	}

	if err := bc.bookService.RateBook(bookId, rateAmount); err != nil {
		errors.SendError(ctx, errors.NewErrRatingBook(err))
		return
	}

	message := "book " + strBookId + " rated with " + strRateAmount

	ctx.JSON(200, gin.H{"message": message,})
}
