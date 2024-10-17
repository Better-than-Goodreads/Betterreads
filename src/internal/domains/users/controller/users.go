package controller

import (
	"net/http"

	"github.com/betterreads/internal/domains/users/models"
	"github.com/betterreads/internal/domains/users/service"
	"github.com/gin-gonic/gin"
)

type UsersController struct {
	us *service.UsersService
}

func NewUsersController(us *service.UsersService) *UsersController {
	return &UsersController{
		us: us,
	}
}

func (u *UsersController) CreateUser(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil { //Esto bindea el JSON a la estructura UserRequest
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} // Faltaria devolverlo lindo	

	userResponse, err := u.us.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": userResponse})
}

func (u *UsersController) GetUsers(c *gin.Context) {
	Users, err := u.us.GetUsers() // Faltaria devolverlo lindo
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"users": Users}) 
	return 
}