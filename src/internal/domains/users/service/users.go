package service

import (
	"errors"
	"github.com/betterreads/internal/domains/users/models"
	"github.com/betterreads/internal/domains/users/repository"
	"github.com/betterreads/internal/domains/users/utils"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	rp repository.UsersDatabase
}

var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrEmailAlreadyTaken    = errors.New("email already taken")
	ErrWrongPassword        = errors.New("wrong password")
)

func NewUsersService(rp repository.UsersDatabase) *UsersService {
	return &UsersService{
		rp: rp,
	}
}

func (u *UsersService) RegisterUser(user *models.UserRequest) (*models.UserResponse, error) {
	stored_user, _ := u.rp.GetUserByUsername(user.Username)
	if stored_user != nil {
		return nil, ErrUsernameAlreadyTaken
	}

	stored_user, _ = u.rp.GetUserByEmail(user.Email)
	if stored_user != nil {
		return nil, ErrEmailAlreadyTaken
	}

	// hashes the password
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, bcrypt.ErrPasswordTooLong // The password max length is 72
	}
	user.Password = hashedPassword

	userRecord, err := u.rp.CreateUser(user)

	if err != nil {
		return nil, err // TODO: return a better error
	}

	UserResponse := utils.MapUserRecordToUserResponse(userRecord)
	return UserResponse, nil
}

func (u *UsersService) LogInUser(user *models.UserLoginRequest) (*models.UserResponse, error) {
	userRecord, err := u.rp.GetUserByUsername(user.Username)
	if userRecord == nil {
		return nil, err
	}

	if !verifyPassword(userRecord.Password, user.Password) {
		return nil, ErrWrongPassword
	}

	UserResponse := utils.MapUserRecordToUserResponse(userRecord)
	return UserResponse, nil
}

func (u *UsersService) GetUsers() ([]*models.UserResponse, error) {
	users, err := u.rp.GetUsers()
	if err != nil {
		return nil, err
	}

	UserResponses, err := utils.MapUsersRecordToUsersResponses(users)

	if err != nil {
		return nil, err
	}

	return UserResponses, nil

}

func (u *UsersService) GetUser(id int) (*models.UserResponse, error) {
	user, err := u.rp.GetUser(id)
	if err != nil {
		return nil, err
	}

	UserResponse := utils.MapUserRecordToUserResponse(user)

	return UserResponse, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
