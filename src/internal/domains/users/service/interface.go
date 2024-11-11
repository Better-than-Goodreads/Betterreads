package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/betterreads/internal/domains/users/models"
	er "github.com/betterreads/internal/pkg/errors"
)

var (
	ErrUsernameTaken = er.ErrorParam{
		Name:   "username",
		Reason: "username already taken",
	}

	ErrEmailTaken = er.ErrorParam{
		Name:   "email",
		Reason: "email already taken",
	}

	ErrUsernameNotFound = errors.New("username doesn't exist")

	ErrWrongPassword = errors.New("wrong password")

	ErrUserNotFound = errors.New("user not found")
)

type UsersService interface {
	RegisterFirstStep(user *models.UserStageRequest) (*models.UserStageResponse, error)
	RegisterSecondStep(user *models.UserAdditionalRequest, id uuid.UUID) (*models.UserResponse, error)
	LogInUser(user *models.UserLoginRequest) (*models.UserResponse, string, error)
	GetUsers() ([]*models.UserResponse, error)
	GetUser(id uuid.UUID) (*models.UserResponse, error)
	PostUserPicture(id uuid.UUID, picture models.UserPictureRequest) error
	GetUserPicture(id uuid.UUID) ([]byte, error)
	SearchUsers(username string, isAuthor bool) ([]*models.UserResponse, error)
    CheckUserExists (id uuid.UUID) bool
}
