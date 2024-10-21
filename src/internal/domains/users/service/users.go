package service

import (
	"errors"

	"github.com/betterreads/internal/domains/users/models"
	rs "github.com/betterreads/internal/domains/users/repository"
    "github.com/betterreads/internal/pkg/auth"
	"github.com/betterreads/internal/domains/users/utils"

)

type UsersService struct {
	rp rs.UsersDatabase
}

var (
	ErrUsernameNotExists    = errors.New("username doesn't exist")
	ErrWrongPassword        = errors.New("wrong password")
)

func NewUsersService(rp rs.UsersDatabase) *UsersService {
	return &UsersService{
		rp: rp,
	}
}


func (u *UsersService) RegisterFirstStep(user *models.UserStageRequest) (*models.UserStageResponse, error) {
    if err := u.rp.CheckUserExists(user); err != nil {
        return nil, err
    }

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return nil, err
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

func (u *UsersService) LogInUser(user *models.UserLoginRequest) (*models.UserResponse, string, error) {
	userRecord, err := u.rp.GetUserByUsername(user.Username)
	if err != nil {
		return nil,"", ErrUsernameNotExists
	}

	if !auth.VerifyPassword(userRecord.Password, user.Password) {
		return nil,"",  ErrWrongPassword
	}

	UserResponse := utils.MapUserRecordToUserResponse(userRecord)
    token, err := auth.GenerateToken(user.Username)
    if err != nil {
        return nil,"", err
    }
	return UserResponse, token, nil
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



func (u *UsersService) GetUser(id string) (*models.UserResponse, error) {
	user, err := u.rp.GetUser(id)
	if err != nil {
		return nil, err
	}

	UserResponse := utils.MapUserRecordToUserResponse(user)

	return UserResponse, nil
}


