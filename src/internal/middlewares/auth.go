package application

import (
	"net/http"
    "strings"
	"github.com/betterreads/internal/pkg/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
    authHeader := c.Request.Header.Get("Authorization")
    if authHeader == ""{
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
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
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
        c.Abort()
        return
    }
    
    c.Set("userId", claims.UserId)
    c.Set("IsAuthor", claims.IsAuthor)

    c.Next()
    
}
