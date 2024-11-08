package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/betterreads/internal/domains/bookshelf/models"
	"github.com/betterreads/internal/domains/bookshelf/service"
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
