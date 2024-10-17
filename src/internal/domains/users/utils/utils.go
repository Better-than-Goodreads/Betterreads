package utils

import (
	"github.com/betterreads/internal/domains/users/models"
)

func MapUsersRecordToUsersResponses(users []models.UserRecord) ([]models.UserResponse, error) {
	userResponses := make([]models.UserResponse, len(users))
	for _ , user:= range(users){
		userResponse := MapUserRecordToUserResponse(user)
		userResponses = append(userResponses, userResponse)
	}
	return userResponses, nil
}

func MapUserRecordToUserResponse(user models.UserRecord) models.UserResponse {
	return models.UserResponse{
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Username: user.Username,
		Location: user.Location,
		Gender: user.Gender,
		Age: user.Age,
		AboutMe: user.AboutMe,
	}
}

func MapUserRequestToUserRecord(user models.UserRequest, id int) models.UserRecord {
	return models.UserRecord{
		Id: id,
		Password: user.Password,
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Username: user.Username,
		Location: user.Location,
		Gender: user.Gender,
		Age: user.Age,
		AboutMe: user.AboutMe,
	}
}	
