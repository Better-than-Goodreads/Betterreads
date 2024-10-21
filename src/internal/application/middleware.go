package application


import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/betterreads/internal/pkg/auth"
)



func authMiddleware(c *gin.Context) {
    // Get the Authorization header
    tokenString := c.Request.Header.Get("Authorization")
    fmt.Printf("token : %s",  tokenString)
    // Check if the token is valid
    if tokenString == ""{
        // Token is missing, abort the request
        c.AbortWithStatus(401)
    }

    claims := &jwt.MapClaims{}
    token , err := jwt.ParseWithClaims(
        tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return auth.JwtSecret, nil
    })

    if err != nil || !token.Valid {
        // Token is invalid, abort the request
        c.JSON(401, gin.H{"error": "Unauthorized request"})
        c.Abort()
    }
    
    username := (*claims)["username"].(string)
    // Attach the username to the context
    c.Set("username", username)

    c.Next()
    
}
