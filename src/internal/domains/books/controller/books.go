package controller

import (
    "errors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/betterreads/internal/domains/books/models"
	"github.com/betterreads/internal/domains/books/service"
	er "github.com/betterreads/internal/pkg/errors"
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
	userId, err := getUserId(ctx)
	if err != nil {
		err := er.NewErrNotLogged()
		ctx.AbortWithError(err.Status, err)
		return
	}

	if !isAuthor {
		err := er.NewErrNotAuthor()
		ctx.AbortWithError(err.Status, err)
		return
	}

	newBookRequest, err := getBookRequest(ctx)
	if err != nil {
		err := er.NewErrParsingRequest(err)
		ctx.AbortWithError(err.Status, err)
		return
	}

	book, err := bc.bookService.PublishBook(newBookRequest, userId)
	if err != nil {
		err := er.NewErrPublishingBook(err)
		ctx.AbortWithError(err.Status, err)
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
		err := er.NewErrAddingReview(err)
		ctx.AbortWithError(err.Status, err)
		return
	}

	book, err := bc.bookService.GetBookInfo(bookUuid, userId)
	if err != nil {
		err := er.NewErrGettingBook(err)
		ctx.AbortWithError(err.Status, err)
		return
	}

	if book == nil {
		err := er.NewErrBookNotFound()
		ctx.AbortWithError(err.Status, err)
		return
	}

	ctx.JSON(http.StatusCreated, book)
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
		err := er.NewErrGettingBooks(err)
		ctx.AbortWithError(err.Status, err)
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
	userId, err := getUserId(ctx)
	if err != nil {
		err := er.NewErrNotLogged()
		ctx.AbortWithError(err.Status, err)
		return
	}

	authorId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := er.NewErrInvalidAuthorId(ctx.Param("id"))
		ctx.AbortWithError(err.Status, err)
		return
	}

	books, err := bc.bookService.GetBooksOfAuthor(authorId, userId)
	if err != nil {
		if errors.Is(err, service.ErrAuthorNotFound) {
			err := er.NewErrAuthorNotFound()
			ctx.AbortWithError(err.Status, err)
			return
		}
		err := er.NewErrGettingBooks(err)
		ctx.AbortWithError(err.Status, err)
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
		err := er.NewErrInvalidBookId(id)
		ctx.AbortWithError(err.Status, err)
		return
	}

	base64Bytes, err := bc.bookService.GetBookPicture(uuid)
	if err != nil {
		err := er.NewErrGettingBook(err)
		ctx.AbortWithError(err.Status, err)
		return
	}

	if base64Bytes == nil {
		err := er.NewErrBookNotFound()
		ctx.AbortWithError(err.Status, err)
		return
	}

	ctx.Data(http.StatusCreated, "image/jpeg", base64Bytes)
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
		err := er.NewErrGettingBooks(err)
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
	userId, err := getUserId(ctx)
	if err != nil {
        err := er.NewErrNotLogged()
        ctx.AbortWithError(err.Status, err)
		return
	}

	var newBookRating models.NewRatingRequest
	if err := ctx.ShouldBindJSON(&newBookRating); err != nil {
        err := er.NewErrParsingRequest(err)
        ctx.AbortWithError(err.Status, err)
		return
	}

	bookId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
        err := er.NewErrInvalidBookId(ctx.Param("id"))
        ctx.AbortWithError(err.Status, err)
		return
	}

	rateAmount := newBookRating.Rating

	rating, err := bc.bookService.RateBook(bookId, userId, rateAmount)
	if err != nil {
        if errors.Is(err, service.ErrBookNotFound) {
            err := er.NewErrBookNotFound()
            ctx.AbortWithError(err.Status, err)
        }
        if errors.Is(err, service.ErrRatingAmount) {
            err := er.NewErrInvalidRatingBook(err)
            ctx.AbortWithError(err.Status, err)
        } else if errors.Is(err, service.ErrRatingAlreadyExists) {
            err := er.NewErrRatingAlreadyExists()
            ctx.AbortWithError(err.Status, err)
        } else {
            err := er.NewErrRating(err)
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
func (bc *BooksController) UpdateRatingOfBook(ctx *gin.Context){
    userId , err := getUserId(ctx)
    if err != nil {
        err := er.NewErrNotLogged()
        ctx.AbortWithError(err.Status, err)
        return
    }

    var newBookRating models.NewRatingRequest
    if err := ctx.ShouldBindJSON(&newBookRating); err != nil {
        err := er.NewErrParsingRequest(err)
        ctx.AbortWithError(err.Status, err)
        return
    }


	bookId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
        err := er.NewErrInvalidBookId(ctx.Param("id"))
        ctx.AbortWithError(err.Status, err)
		return
	}

    err = bc.bookService.UpdateRating(bookId , userId , newBookRating.Rating)
    if err != nil {
        if errors.Is(err, service.ErrRatingNotFound) {
            err := er.NewErrRatingNotFound()
            ctx.AbortWithError(err.Status, err)
        } else if errors.Is(err, service.ErrRatingAmount) {
            err := er.NewErrInvalidRatingBook(err)
            ctx.AbortWithError(err.Status, err)
        } else {
            err := er.NewErrRating(err)
            ctx.AbortWithError(err.Status, err)
        }
    }
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
// func (bc *BooksController) DeleteRating(ctx *gin.Context) {
// 	userId, err := getLoggedUserId(ctx)
// 	if err != nil {
//         err := errors.NewErrNotLogged()
//         ctx.AbortWithError(err.Status, err)
// 		return
// 	}
// 	bookId, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
//         err := errors.NewErrInvalidBookId(ctx.Param("id"))
//         ctx.AbortWithError(err.Status, err)
// 		return
// 	}
//
// 	if err := bc.bookService.DeleteRating(bookId, userId); err != nil {
//         err := errors.NewErrDeletingRating(err)
//         ctx.AbortWithError(err.Status, err)
// 		return
// 	}
//
// 	ctx.JSON(http.StatusNoContent, nil)
// }


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
	userId, err := getUserId(ctx)
	if err != nil {
		err := er.NewErrNotLogged()
		ctx.AbortWithError(err.Status, err)
		return
	}

	bookId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := er.NewErrInvalidBookId(ctx.Param("id"))
		ctx.AbortWithError(err.Status, err)
		return
	}

	var newReview models.NewReviewRequest
	if err := ctx.ShouldBindJSON(&newReview); err != nil {
		err := er.NewErrParsingRequest(err)
		ctx.AbortWithError(err.Status, err)
		return
	}
    fmt.Printf("Review in controller: %v\n", newReview)
    
	if err := bc.bookService.AddReview(bookId, userId, newReview); err != nil {
        if errors.Is(err, service.ErrReviewAlreadyExists) {
            err := er.NewErrReviewAlreadyExists()
            ctx.AbortWithError(err.Status, err)
        } else if errors.Is(err, service.ErrRatingAmount) {
            err := er.NewErrInvalidRatingBook(err)
            ctx.AbortWithError(err.Status, err)
        } else if errors.Is(err, service.ErrBookNotFound) {
            err := er.NewErrBookNotFound()
            ctx.AbortWithError(err.Status, err)
        } else  {
            err := er.NewErrAddingReview(err)
            ctx.AbortWithError(err.Status, err)
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
// @Success 200 {object} []models.Review
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetailsWithParams
// @router /books/{id}/review [get]
func (bc *BooksController) GetBookReviews (ctx *gin.Context) {
    bookId, err := uuid.Parse(ctx.Param("id"))
    if err != nil {
        err := er.NewErrInvalidBookId(ctx.Param("id"))
        ctx.AbortWithError(err.Status, err)
        return
    }

    reviews, err := bc.bookService.GetBookReviews(bookId)
    if err != nil {
        if err == service.ErrBookNotFound {
            err := er.NewErrBookNotFound()
            ctx.AbortWithError(err.Status, err)
        } else {
            err := er.NewErrGettingBookReviews(err)
            ctx.AbortWithError(err.Status, err)
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
func getBookRequest(ctx *gin.Context) (*models.NewBookRequest, error) {
	picture, err := getPicture(ctx)
	if err != nil {
		return nil, err
	}

	data := ctx.PostForm("data")
	var newBookRequest models.NewBookRequest
	if err := json.Unmarshal([]byte(data), &newBookRequest); err != nil {
		return nil, err
	}

	newBookRequest.Picture = picture

	validator := validator.New()
	if err := validator.Struct(newBookRequest); err != nil {
		return nil, err
	}

	return &newBookRequest, nil
}



// Aux
func getPicture(ctx *gin.Context) ([]byte, error) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	picture, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return picture, nil
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
