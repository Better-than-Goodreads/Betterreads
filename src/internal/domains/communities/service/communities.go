package service


import (
	"github.com/betterreads/internal/domains/communities/model"
	"github.com/betterreads/internal/domains/communities/repository"
	"github.com/google/uuid"
	userModel "github.com/betterreads/internal/domains/users/models"
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

func (cs *CommunitiesServiceImpl) JoinCommunity(communityId uuid.UUID, userId uuid.UUID) error {
	if cs.r.CheckIfUserIsInCommunity(communityId, userId) {
		return ErrUserAlreadyInCommunity
	}

	err := cs.r.JoinCommunity(communityId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CommunitiesServiceImpl) GetCommunityUsers(communityId uuid.UUID) ([]*userModel.UserStageResponse, error) {
	users, err := cs.r.GetCommunityUsers(communityId)
	if err != nil {
		return nil, err
	}

	return users, nil
}