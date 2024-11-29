package repository

import (
	"errors"

	"github.com/betterreads/internal/domains/communities/model"
	userModel "github.com/betterreads/internal/domains/users/models"
	"github.com/google/uuid"
)

var (
	ErrCommunityNotFound = errors.New("community not found")
)

type CommunitiesDatabase interface {
	CreateCommunity(community model.NewCommunityRequest, userId uuid.UUID) (*model.CommunityResponse, error)
	GetCommunities(userId uuid.UUID) ([]*model.CommunityResponse, error)
	JoinCommunity(communityId uuid.UUID, userId uuid.UUID) error
	CheckIfUserIsInCommunity(communityId uuid.UUID, userId uuid.UUID) bool
	CheckIFCommunityExists(communityId uuid.UUID) bool
	GetCommunityUsers(communityId uuid.UUID) ([]*userModel.UserStageResponse, error)
	GetCommunityPicture(communityId uuid.UUID) ([]byte, error)
	SearchCommunities(search string, currId uuid.UUID) ([]*model.CommunityResponse, error)
	GetCommunityById(id uuid.UUID, userId uuid.UUID) (*model.CommunityResponse, error)
	GetCommunityPosts(communityId uuid.UUID) ([]*model.CommunityPostResponse, error)
	CreateCommunityPost(communityId uuid.UUID, userId uuid.UUID, content string, title string) error
	LeaveCommunity(communityId uuid.UUID, userId uuid.UUID) error
	CheckIfUserIsCreator(communityId uuid.UUID, userId uuid.UUID) bool
	DeleteCommunity(communityId uuid.UUID) error
}
