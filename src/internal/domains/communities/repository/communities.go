
package repository

import (
	"github.com/betterreads/internal/domains/communities/model"
	"github.com/google/uuid"
)

type CommunitiesDatabase interface {
	CreateCommunity(community model.NewCommunityRequest, userId uuid.UUID) (*model.CommunityResponse, error)
	GetCommunities() ([]*model.CommunityResponse, error)
}
