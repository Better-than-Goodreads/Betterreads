package service

import (
	"errors"

	"github.com/betterreads/internal/domains/communities/model"
	userModel "github.com/betterreads/internal/domains/users/models"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyInCommunity = errors.New("user is already in community")
	ErrUserNotInCommunity     = errors.New("user is not in community")
	ErrCommunityNotFound      = errors.New("community not found")
)

type CommunitiesService interface {
	CreateCommunity(community model.NewCommunityRequest, userId uuid.UUID) (*model.CommunityResponse, error)
	GetCommunities(userId uuid.UUID) ([]*model.CommunityResponse, error)
	JoinCommunity(communityId uuid.UUID, userId uuid.UUID) error
	GetCommunityUsers(communityId uuid.UUID) ([]*userModel.UserStageResponse, error)
	GetCommunityPicture(communityId uuid.UUID) ([]byte, error)
	SearchComunnity(search string, currId uuid.UUID) ([]*model.CommunityResponse, error)
	GetCommunityById(id uuid.UUID, userId uuid.UUID) (*model.CommunityResponse, error)
	GetCommunityPosts(communityId uuid.UUID) ([]*model.CommunityPostResponse, error)
	CreateCommunityPost(communityId uuid.UUID, userId uuid.UUID, content string, title string) error
	LeaveCommunity(communityId uuid.UUID, userId uuid.UUID) error
}
