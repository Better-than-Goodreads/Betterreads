package controller

type CommunitiesController struct {
	communitiesService service.CommunitiesService
}

func NewCommunitiesController(communitiesService service.CommunitiesService) *CommunitiesController {
	return &CommunitiesController{
		communitiesService: communitiesService,
	}
}

func (c *CommunitiesController) CreateCommunity (ctx *gin.Context) {

	userId, errDetail := aux.GetLoggedUserId(ctx)
	if errDetail != nil {
		ctx.AbortWithError(errDetail.Status, errDetail)
		return
	}

	var community model.NewCommunityRequest
	
	if err := ctx.ShouldBindJSON(&community); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCommunity, err := c.communitiesService.CreateCommunity(community)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdCommunity)
}

func (c *CommunitiesController) GetCommunities (ctx *gin.Context) {
	communities, err := c.communitiesService.GetCommunities()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, communities)
}