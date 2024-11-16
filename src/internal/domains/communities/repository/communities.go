
package repository

import (
	"github.com/betterreads/internal/domains/communities/model"
	UUID "github.com/google/uuid"
)

type CommunitiesDatabase interface {
	CreateCommunity(community model.NewCommunityRequest) (UUID.uuid, error)
	GetCommunities() ([]model.CommunityResponse, error)
}
