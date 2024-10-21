package controller

import (
	"net/http"

	"github.com/betterreads/internal/domains/users/models"
	"github.com/betterreads/internal/domains/users/service"
	"github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UsersController struct {
	us *service.UsersService
}

func NewUsersController(us *service.UsersService) *UsersController {
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
		errors.SendError(c, errors.NewErrFetchUsers(err))
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
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		errors.SendError(c, errors.NewErrInvalidUserID(id))
	}

	user, err := u.us.GetUser(uuid)

	if err != nil {
		errors.SendError(c, errors.NewErrUserNotFoundById(err))
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
// @Success 202 {object} models.UserResponse
// @Failure 400 {object} errors.ErrorDetails
// @Failure 400 {object} errors.ErrorDetailsWithParams
// @Failure 404 {object} errors.ErrorDetails
// @Router /users/login [post]
func (u *UsersController) LogIn(c *gin.Context) {
	var user *models.UserLoginRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		errors.SendErrorWithParams(c, errors.NewErrParsingRequest(err))
		return
	}

	userResponse, token, err := u.us.LogInUser(user)

	if err != nil {
		errors.SendError(c, errors.NewErrLogInUser(err))
		return
	}

	c.Header("Authorization", token)

	c.JSON(http.StatusCreated, gin.H{"user": userResponse})
}

func (u *UsersController) Welcome(c *gin.Context) {
	username, _ := c.Get("username")
	msg := "Welcome to BetterReads: " + username.(string)
	c.JSON(http.StatusOK, gin.H{"message": msg})
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
		errors.SendErrorWithParams(c, errors.NewErrParsingRequest(err))
		return
	}

	userStageResponse, err := u.us.RegisterFirstStep(user)

	if err != nil {
		errors.SendError(c, errors.NewErrRegisterUser(err))
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
// @Param id path int true "User register id"
// @Param user body models.UserAdditionalRequest true "User additional request"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} errors.ErrorDetails
// @Failure 400 {object} errors.ErrorDetailsWithParams
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/register/{id}/additional-info [post]
func (u *UsersController) RegisterSecondStep(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		errors.SendError(c, errors.NewErrInvalidRegisterId(id))
		return
	}
	var user *models.UserAdditionalRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		errors.SendErrorWithParams(c, errors.NewErrParsingRequest(err))
		return
	}

	userResponse, err := u.us.RegisterSecondStep(user, uuid)
	if err != nil {
		errors.SendError(c, errors.NewErrRegisterUser(err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": userResponse})
}
