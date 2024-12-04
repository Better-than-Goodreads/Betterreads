package controller

import (
	"errors"
	"net/http"

	"github.com/betterreads/internal/domains/feed/models"
	"github.com/betterreads/internal/domains/feed/service"
	us "github.com/betterreads/internal/domains/users/service"
	aux "github.com/betterreads/internal/pkg/controller"
	er "github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"
)

type FeedController struct {
	fs service.FeedService
}

func NewFeedController(fs service.FeedService) *FeedController {
	return &FeedController{fs: fs}
}

// GetFeed godoc
// @Summary Get feed of an user
// @Description Get feed. The type of posts can be : ["post", "rating"]
// @Tags feed
// @Produce json
// @Error 404 {object} ErrorResponse
// @Error 500 {object} ErrorResponse
// @Success 200 {object} []models.PostDTO
// @Router /feed [get]
func (fc *FeedController) GetFeed(c *gin.Context) {
	userId, errId := aux.GetLoggedUserId(c)
	if errId != nil {
		c.AbortWithError(errId.Status, errId)
	}

	posts, err := fc.fs.GetFeed(userId)

	if err != nil {
		if errors.Is(err, us.ErrUserNotFound) {
			c.AbortWithError(http.StatusNotFound, er.NewErrorDetails("Error when getting feed", err, http.StatusNotFound))
			return
		} else {
			c.AbortWithError(http.StatusInternalServerError, er.NewErrorDetails("Error when getting feed", err, http.StatusInternalServerError))
			return
		}
	}

	res := parsePosts(posts)

	c.JSON(http.StatusOK, res)
}

func parsePosts(posts []models.Post) []models.PostDTO {
	res := make([]models.PostDTO, 0)
	for _, post := range posts {
		var dto models.PostDTO
		if post.Rating != nil {
			dto = models.PostDTO{
				Type: "rating",
				Post: post,
			}
		} else {
			dto = models.PostDTO{
				Type: "post",
				Post: post,
			}
		}
		res = append(res, dto)
	}
	return res
}
