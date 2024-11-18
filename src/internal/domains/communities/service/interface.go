package service

import (
	"github.com/betterreads/internal/domains/communities/model"
	"github.com/google/uuid"
)

type CommunitiesService interface {
	CreateCommunity(community model.NewCommunityRequest, userId uuid.UUID) (*model.CommunityResponse, error)
	GetCommunities() ([]*model.CommunityResponse, error)
}