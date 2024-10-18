package controller

import (
	"github.com/betterreads/internal/domains/users/models"
	"github.com/betterreads/internal/domains/users/service"
	"github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UsersController struct {
	us *service.UsersService
}

func NewUsersController(us *service.UsersService) *UsersController {
	return &UsersController{
		us: us,
	}
}

func (u *UsersController) GetUsers(c *gin.Context) {
	Users, err := u.us.GetUsers()
	if err != nil {
		errors.SendError(c, errors.NewErrFetchUsers(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": Users})
}

func (u *UsersController) GetUser(c *gin.Context) {
	id := c.Param("id")

	id_int, err := strconv.Atoi(id)
	if err != nil {
		errors.SendError(c, errors.NewErrInvalidID(id))
		return
	}

	user, err := u.us.GetUser(id_int)

	if err != nil {
		errors.SendError(c, errors.NewErrUserNotFoundById(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (u *UsersController) LogIn(c *gin.Context) {
	var user *models.UserLoginRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		errors.SendError(c, errors.NewErrParsingRequest(err))
		return
	}

	userResponse, err := u.us.LogInUser(user)

	if err != nil {
		errors.SendError(c, errors.NewErrLogInUser(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": userResponse})
}

func (u *UsersController) Register(c *gin.Context) {
	var user *models.UserRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		errors.SendError(c, errors.NewErrParsingRequest(err))
		return
	}

	userResponse, err := u.us.RegisterUser(user)

	if err != nil {
		errors.SendError(c, errors.NewErrRegisterUser(err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": userResponse})
}
