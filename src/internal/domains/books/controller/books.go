package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/betterreads/internal/domains/books/models"
	"github.com/betterreads/internal/domains/books/service"
	"github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	isAuthor := ctx.GetBool("IsAuthor")
	userId, err := getUserId(ctx)
	if err != nil {
		errors.SendError(ctx, errors.NewErrNotLogged())
		return
	}

	if isAuthor == false {
		errors.SendError(ctx, errors.NewErrNotAuthor())
		return
	}

	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		errors.SendErrorWithParams(ctx, errors.NewErrParsingRequest(err))
		return
	}
	defer file.Close()

	data := ctx.PostForm("data")
	var newBookRequest models.NewBookRequest
	if err := json.Unmarshal([]byte(data), &newBookRequest); err != nil {
		errors.SendErrorWithParams(ctx, errors.NewErrParsingRequest(fmt.Errorf("error parsing JSON request: %w", err)))
		return
	}

	newBookRequest.Picture, err = io.ReadAll(file)
	if err != nil {
		errors.SendErrorWithParams(ctx, errors.NewErrParsingRequest(fmt.Errorf("error reading file: %w", err)))
		return
	}

	validator := validator.New()
	if err := validator.Struct(newBookRequest); err != nil {
		errors.SendErrorWithParams(ctx, errors.NewErrParsingRequest(err))
		return
	}

	book, err := bc.bookService.PublishBook(&newBookRequest, userId)
	if err != nil {
		errors.SendError(ctx, errors.NewErrPublishingBook(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"book": book})
}

// GetBookInfo godoc
// @Summary Get book by id
// @Description Get book id, note that its a UUID
// @Tags books
// @Param id path string true "Book Id"
// @Produce  json
// @Success 200 {object} models.Book
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Router /books/{id}/info [get]
func (bc *BooksController) GetBookInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		errors.SendError(ctx, errors.NewErrInvalidBookId(id))
		return
	}

	book, err := bc.bookService.GetBookInfo(uuid)
	if err != nil {
		errors.SendError(ctx, errors.NewErrGettingBook(err))
		return
	}

	if book == nil {
		errors.SendError(ctx, errors.NewErrBookNotFound())
	}

	ctx.JSON(http.StatusCreated, gin.H{"book": book})
}

// GetBookPicture godoc
// @Summary Get book picture by id
// @Description Get book id, note that its a UUID
// @Tags books
// @Param id path string true "Book Id"
// @Produce  json
// @Success 200 {object} {picture: string}
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Router /books/{id}/picture [get]
func (bc *BooksController) GetBookPicture(ctx *gin.Context) {
	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		errors.SendError(ctx, errors.NewErrInvalidBookId(id))
		return
	}

	book, err := bc.bookService.GetBookPicture(uuid)
	if err != nil {
		errors.SendError(ctx, errors.NewErrGettingBook(err))
		return
	}

	if book == nil {
		errors.SendError(ctx, errors.NewErrBookNotFound())
	}

	ctx.JSON(http.StatusCreated, gin.H{"picture": book})
}

// GetBooksInfo godoc
// @Summary Get all books
// @Description Get all books
// @Tags books
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.Book
// @Failure 500 {object} errors.ErrorDetails
// @Router /books [get]
func (bc *BooksController) GetBooksInfo(ctx *gin.Context) {
	books, err := bc.bookService.GetBooksInfo()
	if err != nil {
		errors.SendError(ctx, errors.NewErrGettingBooks(err))
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"books": books})
}

// RateBook godoc
// @Summary Rate a book
// @Description Rate a book
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "Book Id"
// @Param user body models.NewRatingRequest true "Rating Request"
// @Success 200 {object} string
// @Failure 400 {object} errors.ErrorDetailsWithParams
// @Failure 500 {object} errors.ErrorDetails
// @Router /books/{id}/rating [post]
func (bc *BooksController) RateBook(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		errors.SendError(ctx, errors.NewErrNotLogged())
		return
	}

	var newBookRating models.NewRatingRequest
	if err := ctx.ShouldBindJSON(&newBookRating); err != nil {
		errors.SendErrorWithParams(ctx, errors.NewErrParsingRequest(err))
		return
	}

	bookId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		errors.SendError(ctx, errors.NewErrInvalidBookId(ctx.Param("id")))
		return
	}

	rateAmount := newBookRating.Rating

	if err := bc.bookService.RateBook(bookId, userId, rateAmount); err != nil {
		errors.SendError(ctx, errors.NewErrRatingBook(err))
		return
	}

	ctx.JSON(200, gin.H{"raing": rateAmount})
}

// DeleteRating godoc
// @Summary Delete rating of a book
// @Description Delete rating of a book
// @Tags books
// @Param id path string true "Book Id"
// @Produce  json
// @Success 204 {object} string
// @Failure 400 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /books/{id}/rating [delete]
func (bc *BooksController) DeleteRating(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		errors.SendError(ctx, errors.NewErrNotLogged())
		return
	}
	bookId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		errors.SendError(ctx, errors.NewErrInvalidBookId(ctx.Param("id")))
		return
	}

	if err := bc.bookService.DeleteRating(bookId, userId); err != nil {
		errors.SendError(ctx, errors.NewErrDeletingRating(err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// GetRatingUser godoc
// @Summary Get rating of a book by user
// @Description Get rating of a book by user
// @Tags books
// @Param id path string true "Book Id"
// @Produce  json
// @Success 200 {object} models.RatingResponse
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Router /books/{id}/rating [get]
func (bc *BooksController) GetRatingUser(ctx *gin.Context) {
	bookId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		errors.SendError(ctx, errors.NewErrInvalidBookId(ctx.Param("id")))
		return
	}

	userId, err := getUserId(ctx)
	if err != nil {
		errors.SendError(ctx, errors.NewErrNotLogged())
		return
	}

	ratings, err := bc.bookService.GetRatingUser(bookId, userId)
	if err != nil {
		if err == service.ErrRatingNotFound {
			errors.SendError(ctx, errors.NewErrRatingNotFound())
		} else {
			errors.SendError(ctx, errors.NewErrGettingRating(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ratings": ratings})
}

func getUserId(ctx *gin.Context) (uuid.UUID, error) {
	_userId := ctx.GetString("userId")
	if _userId == "" {
		return uuid.UUID{}, fmt.Errorf("user not logged")
	}
	userId, err := uuid.Parse(_userId)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid user id")
	}
	return userId, nil
}
