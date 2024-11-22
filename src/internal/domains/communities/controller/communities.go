package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/betterreads/internal/domains/communities/model"
	"github.com/betterreads/internal/domains/communities/service"
	aux "github.com/betterreads/internal/pkg/controller"
	er "github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
// @Summary creates a community
// @Description creates a community, the community data should follow the model.NewCommunityRequest in JSON
// @Tags communities
// @Accept  mpfd
// @Produce  json
// @Param data formData string true "Community Data" follows model NewCommunityRequest
// @Param file formData file true "Community Picture"
// @Param community body model.NewCommunityRequest true "Don't need to send this in json, this param is only here to reference NewCommunityRequest, DONT SEND PICTURE in JSON"
// @Success 201 {object} model.CommunityResponse
// @Failure 400 {object} errors.ErrorDetailsWithParams
// @Failure 500 {object} errors.ErrorDetails
// @Router /communities [post]
func (c *CommunitiesController) CreateCommunity(ctx *gin.Context) {
	userId, errDetail := aux.GetLoggedUserId(ctx)
	if errDetail != nil {
		ctx.AbortWithError(errDetail.Status, errDetail)
		return
	}

	newCommunityRequest, errReq := getCommunityRequest(ctx)
	if errReq != nil {
		ctx.AbortWithError(errReq.Status, errReq)
		return
	}

	community, err := c.communitiesService.CreateCommunity(*newCommunityRequest, userId)
	if err != nil {
		errDetail := er.NewErrorDetails("Error when creating community", err, http.StatusInternalServerError)
		ctx.AbortWithError(errDetail.Status, errDetail)
		return
	}

	ctx.JSON(http.StatusCreated, community)
}

// GetCommunities godoc
// @Summary Get all communities
// @Description Get all communities
// @Tags communities
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.CommunityResponse
// @Router /communities [get]
func (c *CommunitiesController) GetCommunities(ctx *gin.Context) {
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
func (c *CommunitiesController) JoinCommunity(ctx *gin.Context) {
	userId, errDetail := aux.GetLoggedUserId(ctx)
	if errDetail != nil {
		ctx.AbortWithError(errDetail.Status, errDetail)
		return
	}

	communityId := ctx.Param("id")
	communityIdParsed, err := uuid.Parse(communityId)
	if err != nil {
		err_detail := fmt.Errorf("Invalid community id: %s", communityId)
		errDetail = er.NewErrorDetails("Error Parsing Community ID", err_detail, http.StatusBadRequest)
		ctx.AbortWithError(errDetail.Status, errDetail)
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
// @Success 200 {array} models.UserStageResponse
// @Router /communities/{id}/users [get]
func (c *CommunitiesController) GetCommunityUsers(ctx *gin.Context) {
	communityId := ctx.Param("id")
	communityIdParsed, err := uuid.Parse(communityId)
	if err != nil {
		err_detail := fmt.Errorf("Invalid community id: %s", communityId)
		errDetail := er.NewErrorDetails("Error Parsing Community ID", err_detail, http.StatusBadRequest)
		ctx.AbortWithError(errDetail.Status, errDetail)
		return
	}

	users, err := c.communitiesService.GetCommunityUsers(communityIdParsed)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// GetPicture godoc
// @Summary Get Community picture
// @Description Get community picture
// @Tags communities
// @Param id path string true "Community id"
// @Produce jpeg
// @Success 200 {file} []byte
// @Failure 400 {object} errors.ErrorDetails
// @Router /communities/{id}/picture [get]
func (c *CommunitiesController) GetCommunityPicture(ctx *gin.Context) {
	communityId := ctx.Param("id")
	communityIdParsed, err := uuid.Parse(communityId)
	if err != nil {
		err_detail := fmt.Errorf("Invalid community id: %s", communityId)
		errDetail := er.NewErrorDetails("Error Parsing Community ID", err_detail, http.StatusBadRequest)
		ctx.AbortWithError(errDetail.Status, errDetail)
		return
	}

	picture, err := c.communitiesService.GetCommunityPicture(communityIdParsed)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(picture) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{})
	}

	ctx.Data(http.StatusOK, "image/jpeg", picture)
}

/*
* getCommunityRequest is a helper function that parses the request body and returns a New
* Community Request struct. It also gets the picture from the request and adds it to the
* NewCommunityRequest struct. It also validates the request and automatically returns an error.
 */
func getCommunityRequest(ctx *gin.Context) (*model.NewCommunityRequest, *er.ErrorDetailsWithParams) {
	picture, err := getPicture(ctx)
	if err != nil {
		return nil, err
	}

	data := ctx.PostForm("data")
	var newCommunityRequest model.NewCommunityRequest
	if err := json.Unmarshal([]byte(data), &newCommunityRequest); err != nil {
		return nil, er.NewErrorDetailsWithParams("Error getting community data", http.StatusBadRequest, err)
	}

	newCommunityRequest.Picture = picture

	validator := validator.New()
	if err := validator.Struct(newCommunityRequest); err != nil {
		return nil, er.NewErrorDetailsWithParams("Error getting community data", http.StatusBadRequest, err)
	}

	return &newCommunityRequest, nil
}

// Aux
func getPicture(ctx *gin.Context) ([]byte, *er.ErrorDetailsWithParams) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		errParam := er.ErrorParam{
			Name:   "picture",
			Reason: "file is required",
		}
		return nil, er.NewErrorDetailsWithParams("Error Creating Community", http.StatusBadRequest, errParam)
	}
	defer file.Close()
	picture, err := io.ReadAll(file)
	if err != nil {
		errParam := er.ErrorParam{
			Name:   "picture",
			Reason: "file is invalid",
		}
		return nil, er.NewErrorDetailsWithParams("Error Creating Community", http.StatusBadRequest, errParam)
	}
	return picture, nil
}
