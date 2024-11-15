package model

import (
	"github.com/google/uuid"
)

type communityResponse struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	OwnerID     uuid.UUID `json:"owner_id" db:"owner_id"`
	//Picture    string    `json:"picture" db:"picture"`
	//Banner   string    `json:"banner" db:"banner"`
	Users 		[]uuid.UUID `json:"users" db:"users"`
	Posts 		[]uuid.UUID `json:"posts" db:"posts"`	
}

type CommunityPost struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CommunityID uuid.UUID `json:"community_id" db:"community_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Content     string    `json:"content" db:"content"`
}

type newCommunityRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	OwnerID     string `json:"owner_id" binding:"required"`
	// Picture     string `json:"picture" binding:"required"`
	// Banner      string `json:"banner" binding:"required"`
}

