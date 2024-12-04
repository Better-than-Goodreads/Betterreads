package controller

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/betterreads/internal/domains/users/models"
	"github.com/betterreads/internal/domains/users/service"
	aux "github.com/betterreads/internal/pkg/controller"
	er "github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UsersController struct {
	us service.UsersService
}

func NewUsersController(us service.UsersService) *UsersController {
	return &UsersController{
		us: us,
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} models.UserResponse
// @Router /users [get]
func (u *UsersController) GetUsers(c *gin.Context) {
	Users, err := u.us.GetUsers()
	if err != nil {
		err := er.NewErrorDetails("Error when getting users", err, http.StatusInternalServerError)
		c.AbortWithError(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": Users})
}

// GetUser godoc
// @Summary Get user by id
// @Description Get user by id
// @Tags users
// @Param id path int true "User id"
// @Produce  json
// @Success 200 {object} models.UserResponse
// @Router /users/{id} [get]
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
func (u *UsersController) GetUser(c *gin.Context) {
	uuid, err := parseUserId(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := u.us.GetUser(uuid)

	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errHttp := er.NewErrorDetails("Error when Getting user", err, http.StatusNotFound)
			c.AbortWithError(errHttp.Status, errHttp)
		} else {
			errHttp := er.NewErrorDetails("Error when Getting user", err, http.StatusInternalServerError)
			c.AbortWithError(errHttp.Status, errHttp)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// LogIn godoc
// @Summary Log in a user
// @Description Log in a user and return a JWT
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.UserLoginRequest true "User login request"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} errors.ErrorDetails
// @Failure 400 {object} errors.ErrorDetailsWithParams
// @Failure 404 {object} errors.ErrorDetails
// @Router /users/login [post]
func (u *UsersController) LogIn(c *gin.Context) {
	var user *models.UserLoginRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		er.AbortWithJsonErorr(c, err)
		return
	}

	userResponse, token, err := u.us.LogInUser(user)

	if err != nil {
		if errors.Is(err, service.ErrUsernameNotFound) {
			errDetails := er.NewErrorDetails("Error when logging in", err, http.StatusNotFound)
			c.AbortWithError(errDetails.Status, errDetails)
		} else if errors.Is(err, service.ErrWrongPassword) {
			errDetails := er.NewErrorDetails("Error when logging in", err, http.StatusUnauthorized)
			c.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when logging in", err, http.StatusInternalServerError)
			c.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  userResponse,
		"token": token,
	})
}

// RegisterBasic godoc
// @Summary Register first step
// @Description Register first step
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.UserStageRequest true "User stage request"
// @Success 201 {object} models.UserStageResponse
// @Failure 400 {object} errors.ErrorDetailsWithParams
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/register/basic [post]
func (u *UsersController) RegisterFirstStep(c *gin.Context) {
	var user *models.UserStageRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		er.AbortWithJsonErorr(c, err)
		return
	}

	userStageResponse, err := u.us.RegisterFirstStep(user)

	if err != nil {
		if errors.Is(err, service.ErrUsernameTaken) || errors.Is(err, service.ErrEmailTaken) {
			errDetails := er.NewErrorDetailsWithParams("Error when registering user", http.StatusBadRequest, err)
			c.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when registering user", err, http.StatusInternalServerError)
			c.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": userStageResponse})
}

// RegisterSecond godoc
// @Summary Register second step
// @Description Register second step
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User register id"
// @Param user body models.UserAdditionalRequest true "User additional request"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} errors.ErrorDetails
// @Failure 400 {object} errors.ErrorDetailsWithParams
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/register/{id}/additional-info [post]
func (u *UsersController) RegisterSecondStep(c *gin.Context) {
	uuid, err := parseUserId(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var user *models.UserAdditionalRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		er.AbortWithJsonErorr(c, err)
		return
	}

	userResponse, err := u.us.RegisterSecondStep(user, uuid)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errDetails := er.NewErrorDetails("Error when registering user", err, http.StatusNotFound)
			c.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Error when registering user", err, http.StatusInternalServerError)
			c.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": userResponse})
}

// PostPicture godoc
// @Summary Post a picture
// @Description Post a picture
// @Tags users
// @Accept  mpfd
// @Produce  json
// @Param file formData file true "User picture"
// @Success 201
// @Failure 400 {object} errors.ErrorDetailsWithParams
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/picture [post]
func (u *UsersController) PostPicture(c *gin.Context) {
	user_id, errId := aux.GetLoggedUserId(c)
	if errId != nil {
		c.AbortWithError(errId.Status, errId)
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		errDetail := fmt.Errorf("Error when parsing picture: %w", err)
		errDetails := er.NewErrorDetails("Failed to post picture", errDetail, http.StatusBadRequest)
		c.AbortWithError(errDetails.Status, errDetails)
		return
	}

	defer file.Close()

	picture, err := io.ReadAll(file)
	if err != nil {
		errDetail := fmt.Errorf("Error when parsing picture: %w", err)
		errDetails := er.NewErrorDetails("Failed to post picture", errDetail, http.StatusInternalServerError)
		c.AbortWithError(errDetails.Status, errDetails)
		return
	}

	request := models.UserPictureRequest{Picture: picture}
	err = u.us.PostUserPicture(user_id, request)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errDetails := er.NewErrorDetails("Failed to post picture", err, http.StatusNotFound)
			c.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Failed to post picture", err, http.StatusInternalServerError)
			c.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// GetPicture godoc
// @Summary Get user picture
// @Description Get user picture
// @Tags users
// @Param id path string true "User id"
// @Produce jpeg
// @Success 200 {file} []byte
// @Failure 400 {object} errors.ErrorDetails
// @Failure 404 {object} errors.ErrorDetails
// @Router /users/{id}/picture [get]
func (u *UsersController) GetPicture(c *gin.Context) {
	user_id, err := parseUserId(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	base64Bytes, err := u.us.GetUserPicture(user_id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errDetails := er.NewErrorDetails("Failed to get user picture", err, http.StatusNotFound)
			c.AbortWithError(errDetails.Status, errDetails)
		} else {
			errDetails := er.NewErrorDetails("Failed to get user picture", err, http.StatusInternalServerError)
			c.AbortWithError(errDetails.Status, errDetails)
		}
		return
	}

	if base64Bytes == nil {
		c.JSON(http.StatusNoContent, gin.H{})
	}

	c.Data(http.StatusOK, "image/jpeg", base64Bytes)
}

// SearchUsers godoc
// @Summary Search users
// @Description Search users
// @Tags users
// @Accept  json
// @Produce  json
// @Param name query string false "User name"
// @Param author query string false "Is author"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/search [get]
func (u *UsersController) SearchUsers(c *gin.Context) {
	name := c.Query("name")
	isAuthor := c.Query("author")

	var users []*models.UserResponse

	var err error
	if isAuthor == "true" {
		users, err = u.us.SearchUsers(name, true)
	} else {
		users, err = u.us.SearchUsers(name, false)
	}

	if err != nil {
		errDetails := er.NewErrorDetails("Error when searching users", err, http.StatusInternalServerError)
		c.AbortWithError(errDetails.Status, errDetails)
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Returns the user id from the context. If the id is not a valid uuid it returns an errorDetails prepared to send.
func parseUserId(ctx *gin.Context) (uuid.UUID, error) {
	id := ctx.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		err_detail := fmt.Errorf("Invalid user id: %s", id)
		return uuid, er.NewErrorDetails("Error when getting user id", err_detail, http.StatusBadRequest)
	}
	return uuid, nil
}
