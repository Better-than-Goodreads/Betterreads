package application

import (
	"time"

	"github.com/gin-gonic/gin"
	"log"
)

// RequestLogger logs the details of each HTTP request
func RequestLogger(c *gin.Context) {
	logger := log.New(log.Writer(), "[GIN] ", log.LstdFlags)

	start := time.Now() // Record the start time

	// Process the request
	c.Next()

	// After request is processed, log the details
	latency := time.Since(start)
	status := c.Writer.Status() // Get the response status code
	clientIP := c.ClientIP()    // Get the client IP
	method := c.Request.Method  // Get the HTTP method
	path := c.Request.URL.Path  // Get the request path

	logger.Printf("%s | %3d | %13v | %15s | %s", method, status, latency, clientIP, path)

}
