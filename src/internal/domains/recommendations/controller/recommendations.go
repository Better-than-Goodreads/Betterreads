package controller

import (
	"errors"
	"net/http"

	_ "github.com/betterreads/internal/domains/books/models" //swagger
	bookService "github.com/betterreads/internal/domains/books/service"
	"github.com/betterreads/internal/domains/recommendations/model"
	"github.com/betterreads/internal/domains/recommendations/service"
	aux "github.com/betterreads/internal/pkg/controller"
	er "github.com/betterreads/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

type RecommenationsController struct {
	rs service.RecommendationsService
}

func NewRecommendationsController(rs service.RecommendationsService) RecommenationsController {
	return RecommenationsController{rs: rs}
}

// GetRecommendations godoc
// @Summary Get recommendations for an user based on his top 3 genres.
// @Description Get recommendations for an user based on his top 3 genres. It gets you 5 books for each genre if available.
// @Tags recommendations
// @Accept  json
// @Produce  json
// @Success 200 {object} []model.RecommendationsByGenre
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /recommendations [get]
func (rc *RecommenationsController) GetRecommendations(c *gin.Context) {
	userId, errId := aux.GetLoggedUserId(c)
	if errId != nil {
		c.AbortWithError(errId.Status, errId)
		return
	}

	booksByTop3Genres, err := rc.rs.GetRecommendations(userId)
	if err != nil {
		if errors.Is(err, bookService.ErrUserNotFound) {
			errDetails := er.NewErrorDetails("Error when getting recommendations", err, http.StatusNotFound)
			c.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrNeedMoreBooksInShelf) {
			errDetails := er.NewErrorDetails("Error when getting recommendations", err, http.StatusBadRequest)
			c.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when getting recommendations", err, http.StatusInternalServerError)
			c.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}
	responses := []model.RecommendationsByGenre{}
	for genre, books := range booksByTop3Genres {
		responses = append(responses, model.RecommendationsByGenre{Genre: genre, Books: books})
	}
	c.JSON(http.StatusOK, responses)
}

// GetMoreRecommendations godoc
// @Summary Get more recommendations for an specific genre.
// @Description Get more recommendations for an specific genre. May want to use it after GetMoreRecommendations, it gets you 20 books if available for the specific genre
// @Tags recommendations
// @Accept  json
// @Produce  json
// @Param genre query string true "Genre"
// @Success 200 {object} []models.Book
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /recommendations/more [get]
func (rc *RecommenationsController) GetMoreRecommendations(c *gin.Context) {
	userId, errId := aux.GetLoggedUserId(c)
	if errId != nil {
		c.AbortWithError(errId.Status, errId)
		return
	}

	genre := c.Query("genre")
	if genre == "" {
		errParam := er.ErrorParam{Name: "genre", Reason: "genre is required"}
		errDetails := er.NewErrorDetailsWithParams("Error when getting recommendations", http.StatusBadRequest, errParam)
		c.AbortWithError(errDetails.Status, errDetails)
		return
	}

	books, err := rc.rs.GetMoreRecommendations(userId, genre)
	if err != nil {
		if errors.Is(err, bookService.ErrGenreNotFound) {
			errDetails := er.NewErrorDetails("Error when getting recommendations", err, http.StatusNotFound)
			c.AbortWithError(errDetails.Status, errDetails)
		}
		if errors.Is(err, bookService.ErrUserNotFound) {
			errDetails := er.NewErrorDetails("Error when getting recommendations", err, http.StatusNotFound)
			c.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrNeedMoreBooksInShelf) {
			errDetails := er.NewErrorDetails("Error when getting recommendations", err, http.StatusBadRequest)
			c.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when getting recommendations", err, http.StatusInternalServerError)
			c.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}

	c.JSON(http.StatusOK, books)
}

// GetFriendsRecommendations godoc
// @Summary Get recommendations for an user based on his friends.
// @Description Get recommendations for an user based on his friends. It gets you all friends books
// @Tags recommendations
// @Produce  json
// @Success 200 {object} []models.Book
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /recommendations/friends [get]
func (rc *RecommenationsController) GetFriendsRecommendations(c *gin.Context) {
	userId, errId := aux.GetLoggedUserId(c)
	if errId != nil {
		c.AbortWithError(errId.Status, errId)
		return
	}
	books, err := rc.rs.GetFriendsRecommendations(userId)
	if err != nil {
		if errors.Is(err, bookService.ErrUserNotFound) {
			errDetails := er.NewErrorDetails("Error when getting recommendations", err, http.StatusNotFound)
			c.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrNeedMoreBooksInShelf) {
			errDetails := er.NewErrorDetails("Error when getting recommendations", err, http.StatusBadRequest)
			c.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when getting recommendations", err, http.StatusInternalServerError)
			c.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}
	c.JSON(http.StatusOK, books)
}
