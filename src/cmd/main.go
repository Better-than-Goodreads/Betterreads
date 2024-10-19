package main

import (
	"github.com/betterreads/internal/application"
    _ "github.com/betterreads/docs"
)

// @title BetterReads API
// @version 1.0
// @description This is a  server for Swagger with Gin.
// @host localhost:8080
// @BasePath /api/v1

func main() {
	r := application.NewRouter(":8080")
	r.Run()
}

