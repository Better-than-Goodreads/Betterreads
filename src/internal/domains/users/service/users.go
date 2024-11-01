package service

import (
	"errors"
    "fmt"
	"github.com/google/uuid"
	"github.com/betterreads/internal/domains/users/models"
	rs "github.com/betterreads/internal/domains/users/repository"
	"github.com/betterreads/internal/domains/users/utils"
	"github.com/betterreads/internal/pkg/auth"
)


type UsersServiceImpl struct {
	rp rs.UsersDatabase
}



func NewUsersServiceImpl(rp rs.UsersDatabase) UsersService{
	return &UsersServiceImpl{
		rp: rp,
	}
}

func (u *UsersServiceImpl) RegisterFirstStep(user *models.UserStageRequest) (*models.UserStageResponse, error) {
    if err := u.rp.CheckUserExistsForRegister(user); err != nil {
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
        return nil, fmt.Errorf("Error when creating stage user: %w",   err)
	}

	UserStageResponse := utils.MapUserStageRecordToUserStageResponse(userRecord)
	return UserStageResponse, nil
}

func (u *UsersServiceImpl)  RegisterSecondStep(user *models.UserAdditionalRequest, id uuid.UUID) (*models.UserResponse, error) {
    UserRecord, err := u.rp.JoinAndCreateUser(user, id)
    if err != nil {
        if errors.Is(err, rs.ErrUserStageNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }

    UserResponse := utils.MapUserRecordToUserResponse(UserRecord)
    return UserResponse, nil
}

func (u *UsersServiceImpl) LogInUser(user *models.UserLoginRequest) (*models.UserResponse, string, error) {
	userRecord, err := u.rp.GetUserByUsername(user.Username)
	if err != nil {
        if errors.Is(err, rs.ErrUserNotFound) {
            return nil,"", ErrUsernameNotFound
        }
		return nil,"", err
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

func (u *UsersServiceImpl) GetUsers() ([]*models.UserResponse, error) {
	users, err := u.rp.GetUsers()
	if err != nil {
		return nil, err
	}

	UserResponses:= utils.MapUsersRecordToUsersResponses(users)

	return UserResponses, nil
}

func (u *UsersServiceImpl) GetUser(id uuid.UUID) (*models.UserResponse, error) {
	user, err := u.rp.GetUser(id)
	if err != nil {
        if errors.Is(err, rs.ErrUserNotFound) {
            return nil, ErrUserNotFound
        }
		return nil, err
	}

	UserResponse := utils.MapUserRecordToUserResponse(user)

	return UserResponse, nil
}

func (u *UsersServiceImpl) PostUserPicture(id uuid.UUID, picture models.UserPictureRequest) error{
    exists := u.rp.CheckUserExists(id)
    if !exists {
        return ErrUserNotFound
    }

    err := u.rp.SaveUserPicture(id, picture.Picture)
    if err != nil {
        if errors.Is(err, rs.ErrUserNotFound) {
            return ErrUserNotFound // to be sure
        }
        return err
    }
    return nil
}

func (u *UsersServiceImpl) GetUserPicture(id uuid.UUID) ([]byte, error) {
    exists := u.rp.CheckUserExists(id)
    if !exists {
        return nil, ErrUserNotFound
    }

    picture, err := u.rp.GetUserPicture(id)
    if err != nil {
        return nil, err
    }
    return picture, nil
}
