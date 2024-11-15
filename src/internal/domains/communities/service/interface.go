package service

type CommunitiesService interface {
	CreateCommunity(community model.NewCommunityRequest) (model.CommunityResponse, error)
}