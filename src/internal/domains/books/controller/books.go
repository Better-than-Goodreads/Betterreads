package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/betterreads/internal/domains/books/models"
	"github.com/betterreads/internal/domains/books/service"
	aux "github.com/betterreads/internal/pkg/controller"
	er "github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BooksController struct {
	bookService service.BooksService
}

func NewBooksController(bookService service.BooksService) *BooksController {
	return &BooksController{bookService: bookService}
}

// PublishBook godoc
// @Summary publish a book
// @Description publishes a book, the book data should follow the models.NewBookRequest in JSON
// @Tags books
// @Accept  mpfd
// @Produce  json
// @Param data formData string true "Book Data" follows model NewBookRequest
// @Param file formData file true "Book Picture"
// @Param book body models.NewBookRequest true "Don't need to send this in json, this param is only here to reference NewBookRequest, DONT SEND PICTURE in JSON"
// @Success 201 {object} models.Book
// @Failure 400 {object} errors.ErrorDetailsWithParams
// @Failure 500 {object} errors.ErrorDetails
// @Router /books [post]
func (bc *BooksController) PublishBook(ctx *gin.Context) {
	isAuthor := ctx.GetBool("IsAuthor")
	userId, errDetail := aux.GetLoggedUserId(ctx)
	if errDetail != nil {
		ctx.AbortWithError(errDetail.Status, errDetail)
		return
	}

	if !isAuthor {
		errDetails := er.NewErrorDetails("Errorr when publishing Book", fmt.Errorf("User is not an author"), http.StatusUnauthorized)
		ctx.AbortWithError(errDetails.Status, errDetails)
		return
	}

	newBookRequest, errReq := getBookRequest(ctx)
	if errReq != nil {
		ctx.AbortWithError(errReq.Status, errReq)
		return
	}

	book, err := bc.bookService.PublishBook(newBookRequest, userId)
	if err != nil {
		if errors.Is(err, service.ErrGenreNotFound) {
			errDetail := er.NewErrorDetailsWithParams("Error when publishing Book", http.StatusBadRequest, err)
			ctx.AbortWithError(errDetail.Status, errDetail)
		} else if errors.Is(err, service.ErrUserNotAuthor) {
			errDetails := er.NewErrorDetails("Error when publishing Book", err, http.StatusUnauthorized)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrAuthorNotFound) {
			errDetail := er.NewErrorDetails("Error when publishing Book", err, http.StatusNotFound)
			ctx.AbortWithError(errDetail.Status, errDetail)
		} else if errors.Is(err, service.ErrGenreRequired) {
			errDetail := er.NewErrorDetailsWithParams("Error when publishing Book", http.StatusBadRequest, err)
			ctx.AbortWithError(errDetail.Status, errDetail)
		} else {
			errDetail := er.NewErrorDetails("Error when publishing Book", err, http.StatusInternalServerError)
			ctx.AbortWithError(errDetail.Status, errDetail)
		}
		return
	}

	ctx.JSON(http.StatusCreated, book)
}

// GetBookInfo godoc
// @Summary Get book by id
// @Description Get book id, note that its a UUID
// @Tags books
// @Param id path string true "Book Id"
// @Produce  json
// @Success 200 {object} models.BookResponseWithReview
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Router /books/{id}/info [get]
func (bc *BooksController) GetBookInfo(ctx *gin.Context) {
	userId := getUserIdIfLogged(ctx)
	bookId := ctx.Param("id")
	bookUuid, err := uuid.Parse(bookId)
	if err != nil {
		errDetails := er.NewErrorDetails("Error when getting Book id", fmt.Errorf("Invalid uuid %s", bookId), http.StatusBadRequest)
		ctx.AbortWithError(errDetails.Status, errDetails)
		return
	}

	book, err := bc.bookService.GetBookInfo(bookUuid, userId)
	if err != nil {
		if errors.Is(err, service.ErrBookNotFound) {
			errDetails := er.NewErrorDetails("Error when getting Book", err, http.StatusNotFound)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when getting Book", err, http.StatusInternalServerError)
			ctx.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}

	ctx.JSON(http.StatusOK, book)
}

// GetBooksByName
// @Summary Get books by name
// @Description Get books by name, if no books found returns an empty array
// @Tags books
// @Param name query string true "Book Name"
// @Produce  json
// @Success 200 {object} []models.BookResponseWithReview
// @Failure 400 {object} errors.ErrorDetails
// @Router /books/info/search [get]
func (bc *BooksController) SearchBooksInfoByName(ctx *gin.Context) {
	userId := getUserIdIfLogged(ctx)
	name := ctx.Query("name")
	books, err := bc.bookService.SearchBooksByName(name, userId)
	if err != nil {
		errDetail := er.NewErrorDetails("Error when searching books", err, http.StatusInternalServerError)
		ctx.AbortWithError(errDetail.Status, errDetail)
		return
	}
	ctx.JSON(http.StatusOK, books)
}

// GetBooksOfAuthor
// @Summary Get books of an auther
// @Description Get the books of an author, if no books found returns an empty array
// @Tags books
// @Param id path string true "Author Id"
// @Produce  json
// @Success 200 {object} []models.BookResponseWithReview
// @Failure 400 {object} errors.ErrorDetails
// @Router /books/author/{id} [get]
func (bc *BooksController) GetBooksOfAuthor(ctx *gin.Context) {
	userId := getUserIdIfLogged(ctx)

	authorId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		errDetail := er.NewErrorDetails("Error when getting Author id", fmt.Errorf("Invalid uuid %s", ctx.Param("id")), http.StatusBadRequest)
		ctx.AbortWithError(errDetail.Status, errDetail)
		return
	}

	books, err := bc.bookService.GetBooksOfAuthor(authorId, userId)
	if err != nil {
		if errors.Is(err, service.ErrAuthorNotFound) {
			errDetails := er.NewErrorDetails("Error when getting books author", err, http.StatusNotFound)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrUserNotAuthor) {
			errDetails := er.NewErrorDetails("Error when getting books author", err, http.StatusUnauthorized)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when getting books author", err, http.StatusInternalServerError)
			ctx.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}
	ctx.JSON(http.StatusOK, books)
}

// GetBookPicture godoc
// @Summary Get book picture by id
// @Description Get book id, note that its a UUID
// @Tags books
// @Param id path string true "Book Id"
// @Produce jpeg
// @Success 200 {file} []byte
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Router /books/{id}/picture [get]
func (bc *BooksController) GetBookPicture(ctx *gin.Context) {
	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		errDetails := er.NewErrorDetails("Error when getting Book id", fmt.Errorf("Invalid id %s", id), http.StatusBadRequest)
		ctx.AbortWithError(errDetails.Status, errDetails)
		return
	}

	base64Bytes, err := bc.bookService.GetBookPicture(uuid)
	if err != nil {
		if errors.Is(err, service.ErrBookNotFound) {
			errDetails := er.NewErrorDetails("Error when getting Book picture", err, http.StatusNotFound)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when getting Book picture", err, http.StatusInternalServerError)
			ctx.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}

	if base64Bytes == nil {
		ctx.JSON(http.StatusNoContent, gin.H{})
	}

	ctx.Data(http.StatusOK, "image/jpeg", base64Bytes)
}

// GetBooksInfo godoc
// @Summary Get all books
// @Description Get all books
// @Tags books
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.BookResponseWithReview
// @Failure 500 {object} errors.ErrorDetails
// @Router /books/info [get]
func (bc *BooksController) GetBooksInfo(ctx *gin.Context) {
	userId := getUserIdIfLogged(ctx)
	books, err := bc.bookService.GetBooksInfo(userId)
	if err != nil {
		err := er.NewErrorDetails("Error when getting books", err, http.StatusInternalServerError)
		ctx.AbortWithError(err.Status, err)
		return
	}
	ctx.JSON(http.StatusAccepted, books)
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

	userId, errDetail := aux.GetLoggedUserId(ctx)
	if errDetail != nil {
		ctx.AbortWithError(errDetail.Status, errDetail)
		return
	}

	var newBookRating models.NewRatingRequest
	if err := ctx.ShouldBindJSON(&newBookRating); err != nil {
		er.AbortWithJsonErorr(ctx, err)
		return
	}

	bookId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		errDetails := er.NewErrorDetails("Error when getting Book id", fmt.Errorf("Invalid uuid %s", ctx.Param("id")), http.StatusBadRequest)
		ctx.AbortWithError(errDetails.Status, errDetails)
		return
	}

	rateAmount := newBookRating.Rating

	rating, err := bc.bookService.RateBook(bookId, userId, rateAmount)
	if err != nil {
		if errors.Is(err, service.ErrBookNotFound) {
			errDetails := er.NewErrorDetails("Error when rating Book", err, http.StatusNotFound)
			ctx.AbortWithError(errDetails.Status, errDetails)
		}
		if errors.Is(err, service.ErrRatingAmount) {
			errDetails := er.NewErrorDetailsWithParams("Error when rating Book", http.StatusBadRequest, err)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrRatingAlreadyExists) {
			errDetails := er.NewErrorDetails("Error when rating Book", err, http.StatusConflict)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrRatingOwnBook) {
			errDetails := er.NewErrorDetails("Error when rating own Book", err, http.StatusForbidden)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else {
			err := er.NewErrorDetails("Error when rating Book", err, http.StatusInternalServerError)
			ctx.AbortWithError(err.Status, err)
		}

		return
	}

	ctx.JSON(200, gin.H{"rating": rating})
}

// UpdateRating godoc
// @Summary Update rating of a book
// @Description Update rating of a book
// @Tags books
// @Param id path string true "Book Id"
// @Produce  json
// @Param user body models.NewRatingRequest true "Rating Request"
// @Success 200 {object} string
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /books/{id}/rating [put]
func (bc *BooksController) UpdateRatingOfBook(ctx *gin.Context) {

	userId, errDetails := aux.GetLoggedUserId(ctx)
	if errDetails != nil {
		ctx.AbortWithError(errDetails.Status, errDetails)
		return
	}

	var newBookRating models.NewRatingRequest
	if err := ctx.ShouldBindJSON(&newBookRating); err != nil {
		er.AbortWithJsonErorr(ctx, err)
		return
	}

	bookId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		errDetails := er.NewErrorDetails("Error when getting Book id", fmt.Errorf("Invalid uuid %s", ctx.Param("id")), http.StatusBadRequest)
		ctx.AbortWithError(errDetails.Status, errDetails)
		return
	}

	err = bc.bookService.UpdateRating(bookId, userId, newBookRating.Rating)
	if err != nil {
		if errors.Is(err, service.ErrRatingNotFound) {
			errDetails := er.NewErrorDetails("Error when updating rating", err, http.StatusNotFound)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrRatingAmount) {
			errDetails := er.NewErrorDetailsWithParams("Error when updating rating", http.StatusBadRequest, errDetails)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when updating rating", err, http.StatusInternalServerError)
			ctx.AbortWithError(errDetails.Status, errDetails)
		}
	}
}

// AddReview godoc
// @Summary Add review to a book
// @Description Add review to a book
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "Book Id"
// @Param user body models.NewReviewRequest true "Review Request"
// @Success 200 {object} string
// @Failure 400 {object} errors.ErrorDetailsWithParams
// @Failure 500 {object} errors.ErrorDetails
// @Router /books/{id}/review [post]
func (bc *BooksController) ReviewBook(ctx *gin.Context) {
	userId, errId := aux.GetLoggedUserId(ctx)
	if errId != nil {
		ctx.AbortWithError(errId.Status, errId)
		return
	}

	bookId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		errDetails := er.NewErrorDetails("Error when getting Book id", fmt.Errorf("Invalid uuid %s", ctx.Param("id")), http.StatusBadRequest)
		ctx.AbortWithError(errDetails.Status, errDetails)
		return
	}

	var newReview models.NewReviewRequest
	if err := ctx.ShouldBindJSON(&newReview); err != nil {
		er.AbortWithJsonErorr(ctx, err)
		return
	}

	if err := bc.bookService.AddReview(bookId, userId, newReview); err != nil {
		if errors.Is(err, service.ErrReviewAlreadyExists) {
			errDetails := er.NewErrorDetails("Error when adding review", err, http.StatusConflict)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrRatingAmount) {
			errDetails := er.NewErrorDetailsWithParams("Error when adding review", http.StatusBadRequest, err)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrBookNotFound) {
			errDetails := er.NewErrorDetails("Error when adding review", err, http.StatusNotFound)
			ctx.AbortWithError(errDetails.Status, errDetails)
        } else if errors.Is(err, service.ErrRatingOwnBook) {
            errDetails := er.NewErrorDetails("Error when rating own Book", err, http.StatusForbidden)
            ctx.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when adding review", err, http.StatusInternalServerError)
			ctx.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}
	ctx.JSON(200, gin.H{"review": newReview.Review})
}

// GetBooksReviews godoc
// @Summary Gets reviews of a book
// @Description Get reviews of a book
// @Tags books
// @Param id path string true "Book Id"
// @Produce  json
// @Success 200 {object} []models.ReviewOfBook
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetailsWithParams
// @router /books/{id}/review [get]
func (bc *BooksController) GetBookReviews(ctx *gin.Context) {
	bookId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		errDetails := er.NewErrorDetails("Error when getting Book id", fmt.Errorf("Invalid uuid %s", ctx.Param("id")), http.StatusBadRequest)
		ctx.AbortWithError(errDetails.Status, errDetails)
		return
	}

	reviews, err := bc.bookService.GetBookReviews(bookId)
	if err != nil {
		if err == service.ErrBookNotFound {
			errDetails := er.NewErrorDetails("Error when getting Book reviews", err, http.StatusNotFound)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when getting Book reviews", err, http.StatusInternalServerError)
			ctx.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

// GetAllReviewsOfUser godoc
// @Summary Gets all reviews of a user
// @Description Get all reviews of a user
// @Tags books
// @Param id path string true "User Id"
// @Produce  json
// @Success 200 {object} []models.ReviewOfUser
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetailsWithParams
// @router /books/user/{id}/reviews [get]
func (bc *BooksController) GetAllReviewsOfUser(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		errDetails := er.NewErrorDetails("Error when getting User id", fmt.Errorf("Invalid uuid %s", ctx.Param("id")), http.StatusBadRequest)
		ctx.AbortWithError(errDetails.Status, errDetails)
		return
	}

	reviews, err := bc.bookService.GetAllReviewsOfUser(userId)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errDetails := er.NewErrorDetails("Error when getting User reviews", err, http.StatusNotFound)
			ctx.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when getting User reviews", err, http.StatusInternalServerError)
			ctx.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

// AUX FUNCTIONS
/*
* getBookRequest is a helper function that parses the request body and returns a New
* Book Request struct. It also gets the picture from the request and adds it to the
* NewBookRequest struct. It also validates the request and automatically sends an error.
 */
func getBookRequest(ctx *gin.Context) (*models.NewBookRequest, *er.ErrorDetailsWithParams) {
	picture, err := getPicture(ctx)
	if err != nil {
		return nil, err
	}

	data := ctx.PostForm("data")
	var newBookRequest models.NewBookRequest
	if err := json.Unmarshal([]byte(data), &newBookRequest); err != nil {
		return nil, er.NewErrorDetailsWithParams("Error getting book data", http.StatusBadRequest, err)
	}

	newBookRequest.Picture = picture

	validator := validator.New()
	if err := validator.Struct(newBookRequest); err != nil {
		return nil, er.NewErrorDetailsWithParams("Error getting book data", http.StatusBadRequest, err)
	}

	return &newBookRequest, nil
}

// Aux
func getPicture(ctx *gin.Context) ([]byte, *er.ErrorDetailsWithParams) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		errParam := er.ErrorParam{
			Name:   "picture",
			Reason: "file is required",
		}
		return nil, er.NewErrorDetailsWithParams("Error Publishing Book", http.StatusBadRequest, errParam)
	}
	defer file.Close()
	picture, err := io.ReadAll(file)
	if err != nil {
		errParam := er.ErrorParam{
			Name:   "picture",
			Reason: "file is invalid",
		}
		return nil, er.NewErrorDetailsWithParams("Error Publishing Book", http.StatusBadRequest, errParam)
	}
	return picture, nil
}

func getUserIdIfLogged(ctx *gin.Context) uuid.UUID {
	_userId := ctx.GetString("userId")
	if _userId == "" {
		return uuid.Nil
	}
	userId, err := uuid.Parse(_userId)
	if err != nil {
		return uuid.Nil
	}
	return userId
}
