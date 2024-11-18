package service


import (
	"github.com/betterreads/internal/domains/communities/model"
	"github.com/betterreads/internal/domains/communities/repository"
	"github.com/google/uuid"
)

type CommunitiesServiceImpl struct {
	r repository.CommunitiesDatabase
	// communitySerice CommunitiesService
}

func NewCommunitiesServiceImpl(r repository.CommunitiesDatabase) CommunitiesService {
	return &CommunitiesServiceImpl{r: r}
}

func (cs *CommunitiesServiceImpl) CreateCommunity(community model.NewCommunityRequest, userId uuid.UUID) (*model.CommunityResponse, error) {
	communityResponse, err := cs.r.CreateCommunity(community, userId)
	if err != nil {
		return nil, err
	}

	return communityResponse, nil
}


func (cs *CommunitiesServiceImpl) GetCommunities() ([]*model.CommunityResponse, error) {
	communities, err := cs.r.GetCommunities()
	if err != nil {
		return nil, err
	}

	return communities, nil
}