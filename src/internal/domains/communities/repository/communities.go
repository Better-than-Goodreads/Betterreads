package repository

import (
	"github.com/betterreads/internal/domains/communities/model"
	userModel "github.com/betterreads/internal/domains/users/models"
	"github.com/google/uuid"
)

type CommunitiesDatabase interface {
	CreateCommunity(community model.NewCommunityRequest, userId uuid.UUID) (*model.CommunityResponse, error)
	GetCommunities(userId uuid.UUID) ([]*model.CommunityResponse, error)
	JoinCommunity(communityId uuid.UUID, userId uuid.UUID) error
	CheckIfUserIsInCommunity(communityId uuid.UUID, userId uuid.UUID) bool
	GetCommunityUsers(communityId uuid.UUID) ([]*userModel.UserStageResponse, error)
	GetCommunityPicture(communityId uuid.UUID) ([]byte, error)
}
