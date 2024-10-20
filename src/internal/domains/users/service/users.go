package service

import (
	"errors"

	"github.com/betterreads/internal/domains/users/models"
	rs "github.com/betterreads/internal/domains/users/repository"
	"github.com/betterreads/internal/domains/users/utils"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	rp rs.UsersDatabase
}

var (
	ErrWrongPassword        = errors.New("wrong password")
)

func NewUsersService(rp rs.UsersDatabase) *UsersService {
	return &UsersService{
		rp: rp,
	}
}


func (u *UsersService) RegisterFirstStep(user *models.UserStageRequest) (*models.UserStageResponse, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, bcrypt.ErrPasswordTooLong
	}

	user.Password = hashedPassword

	userRecord, err := u.rp.CreateStageUser(user)
	if err != nil {
		return nil, err
	}
	UserStageResponse := utils.MapUserStageRecordToUserStageResponse(userRecord)
	return UserStageResponse, nil
}

func (u *UsersService)  RegisterSecondStep(user *models.UserAdditionalRequest) (*models.UserResponse, error) {
    UserRecord, err := u.rp.JoinAndCreateUser(user)
    if err != nil {
        return nil, err
    }

    UserResponse := utils.MapUserRecordToUserResponse(UserRecord)
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
