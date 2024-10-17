package service

import (
	"github.com/betterreads/internal/domains/users/models"
	"github.com/betterreads/internal/domains/users/repository"
	"github.com/betterreads/internal/domains/users/utils"
)

type UsersService struct {
	rp repository.UsersDatabase
}

func NewUsersService(rp repository.UsersDatabase) *UsersService {
	return &UsersService{
		rp: rp,
	}
}

func (u *UsersService) CreateUser(user models.UserRequest) (models.UserResponse, error) {
	userRecord, err := u.rp.CreateUser(user)
	if err != nil {
		return models.UserResponse{}, err //Faltaria mejorar errores
	}

	UserResponse := utils.MapUserRecordToUserResponse(userRecord)
	return UserResponse, nil
}


func (u *UsersService) GetUsers() ([]models.UserResponse, error) {
	users, err := u.rp.GetUsers()
	if err != nil {
		return nil, err
	}

	UserResponses , err:= utils.MapUsersRecordToUsersResponses(users)	

	if err != nil {
		return nil, err
	}

	return UserResponses, nil 

}

func (u *UsersService) GetUser(id int) (models.UserResponse, error) {
	user, err := u.rp.GetUser(id)
	if err != nil {
		return models.UserResponse{}, err
	}

	UserResponse := utils.MapUserRecordToUserResponse(user)	

	return UserResponse, nil 
}

