package controller

import (
    "fmt"
	"errors"
	"io"
	"net/http"

	"github.com/betterreads/internal/domains/users/models"
	"github.com/betterreads/internal/domains/users/service"
	er "github.com/betterreads/internal/pkg/errors"
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
        err := er.NewErrFetchUsers(err)
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
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
        err := er.NewErrInvalidUserID(id)
        c.AbortWithError(err.Status, err)
        return
	}

	user, err := u.us.GetUser(uuid)

	if err != nil {
        err := er.NewErrUserNotFoundById(err)
        c.AbortWithError(err.Status, err)
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
        err := er.NewErrParsingRequest(err)
        c.AbortWithError(err.Status, err)
		return
	}

	userResponse, token, err := u.us.LogInUser(user)

	if err != nil {
        err := er.NewErrLogInUser(err)
        c.AbortWithError(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userResponse,
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
        err := er.NewErrParsingRequest(err)
        c.AbortWithError(err.Status, err)
		return
	}

	userStageResponse, err := u.us.RegisterFirstStep(user)

	if err != nil {
        if errors.Is(err, service.ErrUsernameTaken) || errors.Is(err, service.ErrEmailTaken) {
            err := er.NewErrUserNotUnique(err)
            c.AbortWithError(err.Status, err)
        } else {
            err := er.NewErrRegisterUser(err)
            c.AbortWithError(err.Status, err)
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
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
        err := er.NewErrInvalidRegisterId(id)
        c.AbortWithError(err.Status, err)
		return
	}
	var user *models.UserAdditionalRequest
	if err := c.ShouldBindJSON(&user); err != nil {
        err := er.NewErrParsingRequest(err)
        c.AbortWithError(err.Status, err)
		return
	}

	userResponse, err := u.us.RegisterSecondStep(user, uuid)
	if err != nil {
        err := er.NewErrRegisterUser(err)
        c.AbortWithError(err.Status, err)
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
    user_id , err:= getUserId(c)
    if err != nil {
        err := er.NewErrNotLogged()
        c.AbortWithError(err.Status, err)
    }

    file, _, err := c.Request.FormFile("file")
    if err != nil {
        err := er.NewErrParsingPicture()
        c.AbortWithError(err.Status, err)
        return 
    }

    defer file.Close()

    picture , err := io.ReadAll(file)
    if err != nil {
        err := er.NewErrPostPicture()
        c.AbortWithError(err.Status, err)
        return
    }
    
    request := models.UserPictureRequest{ Picture: picture}
    err = u.us.PostUserPicture(user_id, request)
    if err != nil {
        err := er.NewErrPostPicture() //User must exists because he's logged
        c.AbortWithError(err.Status, err)
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
    id := c.Param("id")
    uuid, err := uuid.Parse(id)
    if err != nil {
        err := er.NewErrInvalidUserID(id)
        c.AbortWithError(err.Status, err)
        return
    }
    base64Bytes, err := u.us.GetUserPicture(uuid)
    if err != nil {
        err := er.NewErrUserNotFoundById(err)
        c.AbortWithError(err.Status, err)
        return
    }

    if base64Bytes == nil {
        err := er.NewErrNoPictureUser()
        c.AbortWithError(err.Status, err)
        return
    }

    c.Data(http.StatusOK, "image/jpeg", base64Bytes)
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
