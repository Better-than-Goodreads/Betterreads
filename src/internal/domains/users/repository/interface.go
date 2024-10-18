package repository

import (
	"errors"

	"github.com/betterreads/internal/domains/users/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UsersDatabase interface {
	CreateUser(user *models.UserRequest) (*models.UserRecord, error)
	GetUser(id int) (*models.UserRecord, error)
	GetUsers() ([]*models.UserRecord, error)
	GetUserByUsername(username string) (*models.UserRecord, error)
	GetUserByEmail(email string) (*models.UserRecord, error)
}
