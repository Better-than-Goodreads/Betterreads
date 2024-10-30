package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/betterreads/internal/domains/users/models"
	rs "github.com/betterreads/internal/domains/users/repository"
	"github.com/betterreads/internal/domains/users/utils"
	"github.com/betterreads/internal/pkg/auth"
	er "github.com/betterreads/internal/pkg/errors"
)

var (
	ErrUsernameTaken = er.ErrorParam{
		Name: "username",
		Reason:  "username already taken",
	}

	ErrEmailTaken = er.ErrorParam {
		Name: "email",
		Reason: "email already taken",
	}

	ErrUserNotFound = errors.New("user not found")
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
		if errors.Is(err, rs.ErrUsernameAlreadyTaken) {
			return nil, ErrUsernameTaken	
		} else if errors.Is(err, rs.ErrEmailAlreadyTaken) {
			return nil, ErrEmailTaken
		}
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

func (u *UsersService)  RegisterSecondStep(user *models.UserAdditionalRequest, id uuid.UUID) (*models.UserResponse, error) {
    UserRecord, err := u.rp.JoinAndCreateUser(user, id)
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

	userResponse := utils.MapUserRecordToUserResponse(userRecord)
    token, err := auth.GenerateToken(userResponse.Id.String(), userResponse.IsAuthor)
    if err != nil {
        return nil,"", err
    }
	return userResponse, token, nil
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



func (u *UsersService) GetUser(id uuid.UUID) (*models.UserResponse, error) {
	user, err := u.rp.GetUser(id)
	if err != nil {
		return nil, rs.ErrUserNotFound
	}

	UserResponse := utils.MapUserRecordToUserResponse(user)

	return UserResponse, nil
}

func (u *UsersService) PostUserPicture(id uuid.UUID, picture models.UserPictureRequest) error{
    err := u.rp.SaveUserPicture(id, picture.Picture)
    if err != nil {
        return err
    }
    return nil
}

func (u *UsersService) GetUserPicture(id uuid.UUID) ([]byte, error) {
    picture, err := u.rp.GetUserPicture(id)
    if err != nil {
		if errors.Is(err, rs.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
        return nil, err
    }
    return picture, nil
}
