package service


type CommunitiesServiceImpl struct {
	r repository.CommunitiesDatabase
	communitySerice service.CommunitiesService
}

func NewCommunitiesServiceImpl(r repository.CommunitiesDatabase, cs service.CommunitiesService) CommunitiesService {
	return &CommunitiesServiceImpl{r: r, communitySerice: cs}
}

func (cs *CommunitiesServiceImpl) CreateCommunity(community model.NewCommunityRequest) (model.CommunityResponse, error) {
	communityId, err := cs.r.CreateCommunity(community)
	if err != nil {
		return model.CommunityResponse{}, err
	}

	communityResponse, err := cs.r.GetCommunityById(communityId)
	if err != nil {
		return model.CommunityResponse{}, err
	}

	return communityResponse, nil
}


func (cs *CommunitiesServiceImpl) GetCommunities() ([]model.CommunityResponse, error) {
	communities, err := cs.r.GetCommunities()
	if err != nil {
		return nil, err
	}

	return communities, nil
}