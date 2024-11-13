package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/betterreads/internal/domains/bookshelf/models"
	"github.com/betterreads/internal/domains/bookshelf/service"
	bookService "github.com/betterreads/internal/domains/books/service"
	aux "github.com/betterreads/internal/pkg/controller"
	er "github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookshelfController struct {
	service service.BookshelfService
}

func NewBookshelfController(service service.BookshelfService) BookshelfController {
	return BookshelfController{service: service}
}

// GetBookShelf godoc
// @Summary Get bookshelf of an user
// @Description Get bookshelf of an user
// @ID get-book
// @Produce  json
// @Tags bookshelf
// @Param userId path string true "User ID"
// @Param type query string true "Shelf Type"
// @Success 200 {object} []models.BookInShelfResponse
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/{userId}/shelf [get]
func (bc *BookshelfController) GetBookShelf(c *gin.Context) {
	userId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errDetails := er.NewErrorDetails("Error when getting shelf", fmt.Errorf("invalid user id"), http.StatusBadRequest)
		c.AbortWithError(errDetails.Status, errDetails)
		return
	}
	shelfType := c.Query("type")
	if shelfType == "" {
		errParam := er.ErrorParam{Name: "type", Reason: "status is required"}
		errDetails := er.NewErrorDetailsWithParams("Error when getting shelf", http.StatusBadRequest, errParam)
		c.AbortWithError(errDetails.Status, errDetails)
		return
	}

	shelf, err := bc.service.GetBookShelf(userId, shelfType)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errDetails := er.NewErrorDetails("Error when getting shelf", err, http.StatusNotFound)
			c.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrInvalidStatusType) {
			errDetails := er.NewErrorDetails("Error when getting shelf", err, http.StatusBadRequest)
			c.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when getting shelf", err, http.StatusInternalServerError)
			c.AbortWithError(errDetails.Status, errDetails)
		}
		return

	}
	c.JSON(200, shelf)
}


// SearchBookShelf godoc
// @Summary Search books in shelf of an user
// @Description Search books in shelf of an user. The search can be filtered by genre, sorted by avg_ratings, total_ratings and date. The direction can be asc or desc.
// @Param status query string true "Shelf Type: all, read, plan-to-read, reading "
// @Param genre query string false "Book Genre"
// @Param sort query string false "Sort by publication_date, total_ratings, avg_rating"
// @Param direction query string false "Sort direction asc or desc"
// @Tags bookshelf
// @Produce  json
// @Success 200 {object} []models.BookInShelfResponse
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/{id}/shelf/search [get]
func (bc *BookshelfController) SearchBookShelf(c *gin.Context) {

    userId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errDetails := er.NewErrorDetails("Error when getting shelf", fmt.Errorf("invalid user id"), http.StatusBadRequest)
		c.AbortWithError(errDetails.Status, errDetails)
		return
	}

    shelfType := c.Query("type")
    if shelfType == "" {
        errParam := er.ErrorParam{Name: "type", Reason: "status is required"}
        errDetails := er.NewErrorDetailsWithParams("Error when getting shelf", http.StatusBadRequest, errParam)
        c.AbortWithError(errDetails.Status, errDetails)
        return
    }

    genre := c.Query("genre")
	sort := c.Query("sort")
	direction := c.Query("direction")
    books, err := bc.service.SearchBookShelf(userId, shelfType, genre, sort, direction)
	if err != nil {
		if errors.Is(err, bookService.ErrInvalidSort) {
			errDetails := er.NewErrorDetailsWithParams(
				"Error when searching books", http.StatusBadRequest, err)
			c.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, bookService.ErrInvalidDirection) {
			errDetails := er.NewErrorDetailsWithParams(
				"Error when searching books", http.StatusBadRequest, err)
			c.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, bookService.ErrDirectionWhenNoSort) {
			errDetails := er.NewErrorDetails("Error when searching books", err, http.StatusBadRequest)
			c.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrGenreNotFound) {
			errDetail := er.NewErrorDetails("Error when searching books in shelf", err, http.StatusBadRequest)
			c.AbortWithError(errDetail.Status, errDetail)
        } else if errors.Is(err, service.ErrInvalidStatusType) {
            errDetails := er.NewErrorDetailsWithParams("Error when searching books in shelf",  http.StatusBadRequest, err)
            c.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetail := er.NewErrorDetails("Error when searching books", err, http.StatusInternalServerError)
			c.AbortWithError(errDetail.Status, errDetail)
		}
		return
	}
	c.JSON(http.StatusOK, books)
}

// AddBookToShelf godoc
// @Summary Add book to shelf
// @Description Add book to shelf
// @ID add-book
// @Accept  json
// @Produce  json
// @Tags bookshelf
// @Param bookShelfEntry body models.BookShelfRequest true "Bookshelf entry"
// @Success 201 {object} string
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetailsWithParams
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/shelf [post]
func (bc *BookshelfController) AddBookToShelf(c *gin.Context) {

	userId, errId := aux.GetLoggedUserId(c)
	if errId != nil {
		c.AbortWithError(errId.Status, errId)
		return
	}

	var req models.BookShelfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		er.AbortWithJsonErorr(c, err)
		return
	}

	err := bc.service.AddBookToShelf(userId, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errDetails := er.NewErrorDetails("Error when adding book to shelf", err, http.StatusNotFound)
			c.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrBookAlreadyInLibrary) {
			errDetails := er.NewErrorDetails("Error when adding book to shelf", err, http.StatusBadRequest)
			c.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrInvalidStatusType) {
			errDetails := er.NewErrorDetails("Error when adding book to shelf", err, http.StatusBadRequest)
			c.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when adding book to shelf", err, http.StatusInternalServerError)
			c.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Book added to shelf"})
}

// EditBookIn godoc
// @Summary Edit book in shelf
// @Description Edit book in shelf
// @ID edit-book
// @Accept  json
// @Produce  json
// @Tags bookshelf
// @Param bookShelfEntry body models.BookShelfRequest true "Bookshelf entry"
// @Success 200 {object} string
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/shelf [put]
func (bc *BookshelfController) EditBookInShelf(c *gin.Context) {
	userId, errId := aux.GetLoggedUserId(c)
	if errId != nil {
		c.AbortWithError(errId.Status, errId)
		return
	}
	var req models.BookShelfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		er.AbortWithJsonErorr(c, err)
		return
	}
	err := bc.service.EditBookInShelf(userId, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errDetails := er.NewErrorDetails("Error when editing book in shelf", err, http.StatusNotFound)
			c.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrBookNotFoundInLibrary) {
			errDetails := er.NewErrorDetails("Error when editing book in shelf", err, http.StatusBadRequest)
			c.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrInvalidStatusType) {
			errDetails := er.NewErrorDetails("Error when editing book in shelf", err, http.StatusBadRequest)
			c.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when editing book in shelf", err, http.StatusInternalServerError)
			c.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book edited in shelf"})

}



// DeleteBookFromShelf godoc
// @Summary Delete book from shelf
// @Description Delete book from shelf
// @ID delete-book
// @Tags bookshelf
// @Param id query string true "Book ID"
// @Produce  json
// @Success 200 {object} string
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/shelf [delete]
func (bc *BookshelfController) DeleteBookFromShelf(c * gin.Context) {
	userId, errId := aux.GetLoggedUserId(c)
	if errId != nil {
		c.AbortWithError(errId.Status, errId)
		return
	}
    
    bookId, err := uuid.Parse(c.Query("id"))
    if err != nil {
        errDetails := er.NewErrorDetails("Error when deleting book from shelf", fmt.Errorf("invalid book id"), http.StatusBadRequest)
        c.AbortWithError(errDetails.Status, errDetails)
        return 
    }

	err = bc.service.DeleteBookFromShelf(userId, bookId)
	if err != nil {if errors.Is(err, service.ErrUserNotFound) {
		errDetails := er.NewErrorDetails("Error when editing book in shelf", err, http.StatusNotFound)
		c.AbortWithError(errDetails.Status, errDetails)
	} else if errors.Is(err, service.ErrBookNotFoundInLibrary) {
		errDetails := er.NewErrorDetails("Error when editing book in shelf", err, http.StatusBadRequest)
		c.AbortWithError(errDetails.Status, errDetails)
	} else {
		errDetails := er.NewErrorDetails("Error when editing book in shelf", err, http.StatusInternalServerError)
		c.AbortWithError(errDetails.Status, errDetails)
	}
	return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted from shelf"})
}
