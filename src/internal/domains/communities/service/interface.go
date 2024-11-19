package service

import (
	"github.com/betterreads/internal/domains/communities/model"
	"github.com/google/uuid"
	userModel "github.com/betterreads/internal/domains/users/models"
	"errors"
)

var (
	ErrUserAlreadyInCommunity  = errors.New("user is already in community")

)

type CommunitiesService interface {
	CreateCommunity(community model.NewCommunityRequest, userId uuid.UUID) (*model.CommunityResponse, error)
	GetCommunities(userId uuid.UUID) ([]*model.CommunityResponse, error)
	JoinCommunity(communityId uuid.UUID, userId uuid.UUID) error
	GetCommunityUsers(communityId uuid.UUID) ([]*userModel.UserStageResponse, error)
}