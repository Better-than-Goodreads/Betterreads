package application

import (
	"fmt"
	"github.com/betterreads/internal/pkg/auth"
	er "github.com/betterreads/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		er.SendError(c, er.NewErrorDetails("Not authorized", fmt.Errorf("Invalid Log in attempt"), http.StatusUnauthorized))
		c.Abort()
		return
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		c.Abort()
		return
	}

	tokenString := bearerToken[1]
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		er.SendError(c, er.NewErrorDetails("Not authorized", fmt.Errorf("Invalid Log in attempt"), http.StatusUnauthorized))
		c.Abort()
		return
	}

	c.Set("userId", claims.UserId)
	c.Set("IsAuthor", claims.IsAuthor)

	c.Next()

}

// Middleware to authenticate users in public routes. If it finds a token it sets it in the context. Else it leaves it empty.
func AuthPublicMiddleware(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.Next()
		return
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		c.Abort()
		return
	}
	tokenString := bearerToken[1]
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		c.Abort()
		return
	}

	c.Set("userId", claims.UserId)
	c.Set("IsAuthor", claims.IsAuthor)
	c.Next()

}
