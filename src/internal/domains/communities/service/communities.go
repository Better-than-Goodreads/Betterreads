package service

import (
	"github.com/betterreads/internal/domains/communities/model"
	"github.com/betterreads/internal/domains/communities/repository"
	userModel "github.com/betterreads/internal/domains/users/models"
	"github.com/google/uuid"
)

type CommunitiesServiceImpl struct {
	r repository.CommunitiesDatabase
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

func (cs *CommunitiesServiceImpl) GetCommunities(userId uuid.UUID) ([]*model.CommunityResponse, error) {
	communities, err := cs.r.GetCommunities(userId)
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

func (cs *CommunitiesServiceImpl) GetCommunityPicture(communityId uuid.UUID) ([]byte, error) {
	picture, err := cs.r.GetCommunityPicture(communityId)
	if err != nil {
		return nil, err
	}
	return picture, nil
}

func (cs *CommunitiesServiceImpl) SearchComunnity(search string, currId uuid.UUID) ([]*model.CommunityResponse, error) {
	communities, err := cs.r.SearchCommunities(search, currId)

	if err != nil {
		return nil, err
	}

	if len(communities) == 0 {
		return []*model.CommunityResponse{}, nil
	}

	return communities, nil
}

func (cs *CommunitiesServiceImpl) GetCommunityById(id uuid.UUID, userId uuid.UUID) (*model.CommunityResponse, error) {
	community, err := cs.r.GetCommunityById(id, userId)
	if err != nil {
		if err == repository.ErrCommunityNotFound {
			return nil, ErrCommunityNotFound
		}
		return nil, err
	}
	return community, nil
}

func (cs *CommunitiesServiceImpl) GetCommunityPosts(communityId uuid.UUID) ([]*model.CommunityPostResponse, error) {
	exists := cs.r.CheckIFCommunityExists(communityId)
	if !exists {
		return nil, ErrCommunityNotFound
	}

	posts, err := cs.r.GetCommunityPosts(communityId)
	if err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return []*model.CommunityPostResponse{}, nil
	}

	return posts, nil
}

func (cs *CommunitiesServiceImpl) CreateCommunityPost(communityId uuid.UUID, userId uuid.UUID, content string, title string) error {
	userInCommunity := cs.r.CheckIfUserIsInCommunity(communityId, userId)
	if !userInCommunity {
		return ErrUserNotInCommunity
	}

	err := cs.r.CreateCommunityPost(communityId, userId, content, title)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CommunitiesServiceImpl) LeaveCommunity(communityId uuid.UUID, userId uuid.UUID) error {
	userInCommunity := cs.r.CheckIfUserIsInCommunity(communityId, userId)
	if !userInCommunity {
		return ErrUserNotInCommunity
	}

	err := cs.r.LeaveCommunity(communityId, userId)
	if err != nil {
		return err
	}
	return nil
}
