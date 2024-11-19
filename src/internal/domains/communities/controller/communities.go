package controller

import(
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/betterreads/internal/domains/communities/service"
	"github.com/betterreads/internal/domains/communities/model"
	aux "github.com/betterreads/internal/pkg/controller"
	"github.com/google/uuid"
)

type CommunitiesController struct {
	communitiesService service.CommunitiesService
}

func NewCommunitiesController(communitiesService service.CommunitiesService) *CommunitiesController {
	return &CommunitiesController{
		communitiesService: communitiesService,
	}
}

// CreateCommunity godoc
// @Summary Create a new community
// @Description Create a new community
// @Tags communities
// @Accept json
// @Produce json
// @Param community body NewCommunityRequest true "Community object that needs to be created"
// @Security ApiKeyAuth
// @Success 201 {object} CommunityResponse
// @Router /communities [post]
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

	createdCommunity, err := c.communitiesService.CreateCommunity(community, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdCommunity)
}

// GetCommunities godoc
// @Summary Get all communities
// @Description Get all communities
// @Tags communities
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} CommunityResponse
// @Router /communities [get]

func (c *CommunitiesController) GetCommunities (ctx *gin.Context) {
	userId, errDetail := aux.GetLoggedUserId(ctx)
	if errDetail != nil {
		ctx.AbortWithError(errDetail.Status, errDetail)
		return
	}
	communities, err := c.communitiesService.GetCommunities(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, communities)
}

// JoinCommunity godoc
// @Summary Join a community
// @Description Join a community
// @Tags communities
// @Accept json
// @Produce json
// @Param id path string true "Community ID"
// @Security ApiKeyAuth
// @Success 200 {string} string
// @Router /communities/{id}/join [post]

func (c *CommunitiesController) JoinCommunity (ctx *gin.Context) {
	userId, errDetail := aux.GetLoggedUserId(ctx)
	if errDetail != nil {
		ctx.AbortWithError(errDetail.Status, errDetail)
		return
	}

	communityId := ctx.Param("id")
	communityIdParsed, err := uuid.Parse(communityId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err2 := c.communitiesService.JoinCommunity(communityIdParsed, userId)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User joined community"})
}

// GetCommunityUsers godoc
// @Summary Get all users in a community
// @Description Get all users in a community
// @Tags communities
// @Accept json
// @Produce json
// @Param id path string true "Community ID"
// @Security ApiKeyAuth
// @Success 200 {array} UserStageResponse
// @Router /communities/{id}/users [get]

func (c *CommunitiesController) GetCommunityUsers (ctx *gin.Context) {
	communityId := ctx.Param("id")
	communityIdParsed, err := uuid.Parse(communityId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users, err := c.communitiesService.GetCommunityUsers(communityIdParsed)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}