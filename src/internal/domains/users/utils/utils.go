package utils

import (
	"github.com/betterreads/internal/domains/users/models"
)

func MapUsersRecordToUsersResponses(users []*models.UserRecord) ([]*models.UserResponse, error) {
	userResponses := make([]*models.UserResponse, 0, len(users))
	for _, user := range users {
		userResponse := MapUserRecordToUserResponse(user)
		userResponses = append(userResponses, userResponse)
	}
	return userResponses, nil
}

func MapUserRecordToUserResponse(user *models.UserRecord) *models.UserResponse {
	return &models.UserResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Location:  user.Location,
		Gender:    user.Gender,
		Age:       user.Age,
		AboutMe:   user.AboutMe,
	}
}

func MapUserRequestToUserRecord(user *models.UserRequest, id int) *models.UserRecord {
	return &models.UserRecord{
		Id:        id,
		Password:  user.Password,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Location:  user.Location,
		Gender:    user.Gender,
		Age:       user.Age,
		AboutMe:   user.AboutMe,
	}
}


func MapUserStageRecordToUserStageResponse (user *models.UserStageRecord) *models.UserStageResponse {
    return &models.UserStageResponse{
        Email: user.Email,
        Username: user.Username,
        First_name: user.FirstName,
        Last_name: user.LastName,
        Token: user.Token,
    }
}


func MapUserStageRequestToUserStageRecord (user *models.UserStageRequest, token string) *models.UserStageRecord {
    return &models.UserStageRecord{
        Email: user.Email,
        Username: user.Username,
        Password: user.Password,
        FirstName: user.FirstName,
        LastName: user.LastName,
        Token: token,
    }
}

func CombineUser(userPrimary *models.UserStageRecord, userSecondary *models.UserAdditionalRequest) *models.UserRecord {
    return &models.UserRecord{
        Email: userPrimary.Email,
        Password: userPrimary.Password,
        FirstName: userPrimary.FirstName,
        LastName: userPrimary.LastName,
        Username: userPrimary.Username,
        Location: userSecondary.Location,
        Gender: userSecondary.Gender,
        Age: userSecondary.Age,
        AboutMe: userPrimary.Username,
    }
}
