package model

import (
	"github.com/google/uuid"
)

type CommunityResponse struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	OwnerID     uuid.UUID `json:"owner_id" db:"owner_id"`
	Joined	  bool      `json:"joined" db:"joined"`
	//Picture    string    `json:"picture" db:"picture"`
	//Banner   string    `json:"banner" db:"banner"`
}

type NewCommunityRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`

	// Picture     string `json:"picture" binding:"required"`
	// Banner      string `json:"banner" binding:"required"`
}

