package controller

import (
	"errors"
	"fmt"
	_ "github.com/betterreads/internal/domains/users/models"
	"github.com/betterreads/internal/domains/friends/service"
	usersService "github.com/betterreads/internal/domains/users/service"
	aux "github.com/betterreads/internal/pkg/controller"
	er "github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type FriendsController struct {
	FriendsService service.FriendsService
}

func NewFriendsController(fs service.FriendsService) FriendsController {
	return FriendsController{FriendsService: fs}
}

// GetFriends godoc
// @Summary Get Friends
// @Description Get Friends of user logged in
// @Tags Friends
// @Param id path string true "User ID"
// @Produce json
// @Success 200 {array} []models.UserResponse
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/{id}/friends [get]
func (fc *FriendsController) GetFriends(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		err := fmt.Errorf("Invalid User ID")
		errorDetails := er.NewErrorDetails("Error When getting Friends", err, http.StatusBadRequest)
		ctx.AbortWithError(http.StatusBadRequest, errorDetails)
		return
	}

	friends, err := fc.FriendsService.GetFriends(id)
	if err != nil {
		if errors.Is(err, usersService.ErrUserNotFound) {
			errorDetails := er.NewErrorDetails("Error When getting Friends", err, http.StatusNotFound)
			ctx.AbortWithError(http.StatusNotFound, errorDetails)
		} else {
			errorDetails := er.NewErrorDetails("Error When getting Friends", err, http.StatusInternalServerError)
			ctx.AbortWithError(http.StatusInternalServerError, errorDetails)
		}
		return
	}
	ctx.JSON(http.StatusOK, friends)
}

// AddFriend godoc
// @Summary Add a friend
// @Description Add Friend to user logged in
// @Tags Friends
// @Produce json
// @Param Id query string true "Friend ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} errors.ErrorDetails
// @Failure 409 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/friends [post]
func (fc *FriendsController) AddFriend(ctx *gin.Context) {
	senderId, errId := aux.GetLoggedUserId(ctx)
	if errId != nil {
		ctx.AbortWithError(http.StatusUnauthorized, errId)
		return
	}

	recipientId, err := uuid.Parse(ctx.Query("Id"))
	if err != nil {
		errorDetails := er.NewErrorDetails("Error When adding friend", err, http.StatusBadRequest)
		ctx.AbortWithError(http.StatusBadRequest, errorDetails)
		return
	}

	err = fc.FriendsService.AddFriend(senderId, recipientId)
	if err != nil {
		if errors.Is(err, usersService.ErrUserNotFound) {
			errorDetails := er.NewErrorDetails("Error When adding friend", err, http.StatusNotFound)
			ctx.AbortWithError(http.StatusNotFound, errorDetails)
		} else if errors.Is(err, service.ErrSameUser) {
			errorDetails := er.NewErrorDetails("Error When adding friend", err, http.StatusForbidden)
			ctx.AbortWithError(http.StatusBadRequest, errorDetails)
		} else if errors.Is(err, service.ErrUserFriendNotFound) {
			errorDetails := er.NewErrorDetails("Error When adding friend", err, http.StatusNotFound)
			ctx.AbortWithError(http.StatusNotFound, errorDetails)
		} else if errors.Is(err, service.ErrFriendRequestExists) {
			errorDetails := er.NewErrorDetails("Error When adding friend", err, http.StatusConflict)
			ctx.AbortWithError(http.StatusConflict, errorDetails)
		} else if errors.Is(err, service.ErrAlreadyFriends) {
			errorDetails := er.NewErrorDetails("Error When adding friend", err, http.StatusConflict)
			ctx.AbortWithError(http.StatusConflict, errorDetails)
		} else {
			errorDetails := er.NewErrorDetails("Error When adding friend", err, http.StatusInternalServerError)
			ctx.AbortWithError(http.StatusInternalServerError, errorDetails)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Friend request sent"})
}

// AcceptFriendRequest godoc
// @Summary Accept a friend request
// @Description Accept Friend Request from user logged in
// @Tags Friends
// @Produce json
// @Param Id query string true "Friend ID"
// @Success 200 {object} string
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/friends/requests [post]
func (fc *FriendsController) AcceptFriendRequest(ctx *gin.Context) {
	recipientId, errId := aux.GetLoggedUserId(ctx)
	if errId != nil {
		ctx.AbortWithError(http.StatusUnauthorized, errId)
		return
	}
	senderId, err := uuid.Parse(ctx.Query("Id"))
	if err != nil {
		errorDetails := er.NewErrorDetails("Error When accepting friend", err, http.StatusBadRequest)
		ctx.AbortWithError(http.StatusBadRequest, errorDetails)
		return
	}
	err = fc.FriendsService.AcceptFriendRequest(recipientId, senderId)
	if err != nil {
		if errors.Is(err, service.ErrRequestNotFound) {
			errorDetails := er.NewErrorDetails("Error When accepting friend", err, http.StatusNotFound)
			ctx.AbortWithError(http.StatusNotFound, errorDetails)
		} else {
			errorDetails := er.NewErrorDetails("Error When accepting friend", err, http.StatusInternalServerError)
			ctx.AbortWithError(http.StatusInternalServerError, errorDetails)
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Friend request accepted"})
}

// RejectFriendRequest godoc
// @Summary Reject a friend request
// @Description Reject Friend Request from
// @Tags Friends
// @Produce json
// @Param Id query string true "Friend ID"
// @Success 200 {object} string
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/friends/requests [delete]
func (fc FriendsController) RejectFriendRequest(ctx *gin.Context) {
	recipientId, errId := aux.GetLoggedUserId(ctx)
	if errId != nil {
		ctx.AbortWithError(http.StatusUnauthorized, errId)
		return
	}
	senderId, err := uuid.Parse(ctx.Query("Id"))
	if err != nil {
		errorDetails := er.NewErrorDetails("Error When declining friend", err, http.StatusBadRequest)
		ctx.AbortWithError(http.StatusBadRequest, errorDetails)
		return
	}
	err = fc.FriendsService.RejectFriendRequest(recipientId, senderId)
	if err != nil {
		if errors.Is(err, service.ErrRequestNotFound) {
			errorDetails := er.NewErrorDetails("Error When declining friend", err, http.StatusNotFound)
			ctx.AbortWithError(http.StatusNotFound, errorDetails)
		} else {
			errorDetails := er.NewErrorDetails("Error When declining friend", err, http.StatusInternalServerError)
			ctx.AbortWithError(http.StatusInternalServerError, errorDetails)
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Friend request declined"})
}

// GetFriends
// @Summary Get Friends Request Sent
// @Description Get Friends Request Sent
// @Tags Friends
// @Produce json
// @Success 200 {array} []models.UserResponse
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/friends/requests/sent [get]
func (fc FriendsController) GetFriendsRequestSent(ctx *gin.Context) {
	id, errId := aux.GetLoggedUserId(ctx)
	if errId != nil {
		ctx.AbortWithError(http.StatusUnauthorized, errId)
		return
	}
	friends, err := fc.FriendsService.GetFriendRequestsSent(id)
	if err != nil {
		if errors.Is(err, usersService.ErrUserNotFound) {
			errorDetails := er.NewErrorDetails("Error When getting Friends Request Sent", err, http.StatusNotFound)
			ctx.AbortWithError(http.StatusNotFound, errorDetails)
		} else {
			errorDetails := er.NewErrorDetails("Error When getting Friends Request Sent", err, http.StatusInternalServerError)
			ctx.AbortWithError(http.StatusInternalServerError, errorDetails)
		}
		return
	}
	ctx.JSON(http.StatusOK, friends)
}

// GetFriends
// @Summary Get Friends Request Received
// @Description Get Friends Request Received
// @Tags Friends
// @Produce json
// @Success 200 {object} []models.UserResponse
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/friends/requests/received [get]
func (fc FriendsController) GetFriendRequestsReceived(ctx *gin.Context) {
	id, errId := aux.GetLoggedUserId(ctx)
	if errId != nil {
		ctx.AbortWithError(http.StatusUnauthorized, errId)
		return
	}
	friends, err := fc.FriendsService.GetFriendRequestsReceived(id)
	if err != nil {
		if errors.Is(err, usersService.ErrUserNotFound) {
			errorDetails := er.NewErrorDetails("Error When getting Friends Request Received", err, http.StatusNotFound)
			ctx.AbortWithError(http.StatusNotFound, errorDetails)
		} else {
			errorDetails := er.NewErrorDetails("Error When getting Friends Request Received", err, http.StatusInternalServerError)
			ctx.AbortWithError(http.StatusInternalServerError, errorDetails)
		}
		return
	}
	ctx.JSON(http.StatusOK, friends)
}

// DeleteFriend godoc
// @Summary Delete a friend
// @Description Delete Friend from user logged in
// @Tags Friends
// @Produce json
// @Param Id query string true "Friend ID"
// @Success 200
// @Failure 404 {object} errors.ErrorDetails
// @Failure 500 {object} errors.ErrorDetails
// @Router /users/friends [delete]
func (fc FriendsController) DeleteFriend(ctx *gin.Context){
    id, errId := aux.GetLoggedUserId(ctx)
    if errId !=nil {
        ctx.AbortWithError(http.StatusUnauthorized, errId)
        return
    }

	friendId, err := uuid.Parse(ctx.Query("Id"))
	if err != nil {
		errorDetails := er.NewErrorDetails("Error When delete friend", err, http.StatusBadRequest)
		ctx.AbortWithError(http.StatusBadRequest, errorDetails)
		return
	}

    if err := fc.FriendsService.DeleteFriend(id, friendId); err != nil {
        if errors.Is(err, service.ErrFriendShipNotFound) {
            errorDetails := er.NewErrorDetails("Error When deleting friend", err, http.StatusNotFound)
            ctx.AbortWithError(http.StatusNotFound, errorDetails)
        } else {
            errorDetails := er.NewErrorDetails("Error When deleting friend", err, http.StatusInternalServerError)
            ctx.AbortWithError(http.StatusInternalServerError, errorDetails)
        }
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Friend deleted"})
}
