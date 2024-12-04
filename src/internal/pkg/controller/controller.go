package controller

import (
	"fmt"
	"net/http"

	er "github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetLoggedUserId(ctx *gin.Context) (uuid.UUID, *er.ErrorDetails) {
	_userId := ctx.GetString("userId")
	if _userId == "" {
		err := er.NewErrorDetails("Error when getting User id", fmt.Errorf("invalid user id: %s", _userId), http.StatusUnauthorized)
		return uuid.UUID{}, err
	}
	userId, err := uuid.Parse(_userId)
	if err != nil {
		err := er.NewErrorDetails("Error when getting User id", fmt.Errorf("invalid user id: %s", _userId), http.StatusUnauthorized)
		return uuid.UUID{}, err
	}
	return userId, nil
}

func GetUserIdIfLogged(ctx *gin.Context) uuid.UUID {
	_userId := ctx.GetString("userId")
	if _userId == "" {
		return uuid.Nil
	}
	userId, err := uuid.Parse(_userId)
	if err != nil {
		return uuid.Nil
	}
	return userId
}
