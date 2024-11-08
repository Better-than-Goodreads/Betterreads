package utils

import (
	"github.com/betterreads/internal/domains/users/models"
)

func MapUsersRecordToUsersResponses(users []*models.UserRecord) []*models.UserResponse {
	userResponses := make([]*models.UserResponse, 0, len(users))
	for _, user := range users {
		userResponse := MapUserRecordToUserResponse(user)
		userResponses = append(userResponses, userResponse)
	}
	return userResponses
}

func MapUserRecordToUserResponse(user *models.UserRecord) *models.UserResponse {
	return &models.UserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Location:  user.Location,
		Gender:    user.Gender,
		Id:        user.Id,
		Age:       user.Age,
		AboutMe:   user.AboutMe,
		IsAuthor:  user.IsAuthor,
	}
}

func MapUserStageRecordToUserStageResponse(user *models.UserStageRecord) *models.UserStageResponse {
	return &models.UserStageResponse{
		Email:      user.Email,
		Username:   user.Username,
		First_name: user.FirstName,
		Last_name:  user.LastName,
		Id:         user.Id,
		IsAuthor:   user.IsAuthor,
	}
}

func MapUserStageRequestToUserStageRecord(user *models.UserStageRequest) *models.UserStageRecord {
	return &models.UserStageRecord{
		Email:     user.Email,
		Username:  user.Username,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsAuthor:  user.IsAuthor,
	}
}

func CombineUser(userPrimary *models.UserStageRecord, userSecondary *models.UserAdditionalRequest) *models.UserRecord {
	return &models.UserRecord{
		Email:     userPrimary.Email,
		Password:  userPrimary.Password,
		FirstName: userPrimary.FirstName,
		LastName:  userPrimary.LastName,
		Username:  userPrimary.Username,
		Location:  userSecondary.Location,
		Gender:    userSecondary.Gender,
		Age:       userSecondary.Age,
		AboutMe:   userPrimary.Username,
		IsAuthor:  userPrimary.IsAuthor,
	}
}
