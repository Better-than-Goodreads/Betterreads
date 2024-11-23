package model

import (
	"github.com/google/uuid"
)

type CommunityResponse struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	OwnerID     uuid.UUID `json:"owner_id" db:"owner_id"`
	Joined      bool      `json:"joined" db:"joined"`
}

type NewCommunityRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Picture     []byte `json:"picture" `
}

type NewCommunityPostRequest struct {
	Content string `json:"content" binding:"required"`
}
