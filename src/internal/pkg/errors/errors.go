package errors

import (
	gin "github.com/gin-gonic/gin"
)

// Follows RFC 7807: https://datatracker.ietf.org/doc/html/rfc7807
type ErrorDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
	Status   int    `json:"status"`
}

func (e *ErrorDetails) Error() string {
	return e.Detail
}

func NewErrorDetails(title string, detail string, status int) *ErrorDetails {
	return &ErrorDetails{
		Type:   "about:blank",
		Title:  title,
		Detail: detail,
		Status: status,
	}
}

// TODO: Implement errors with validations for the json requests
// With the error_params:
//      - name: "email"
//      - reason: "Email is required"

func SendError(c *gin.Context, err *ErrorDetails) {
	err.Instance = c.Request.RequestURI
	c.JSON(err.Status, err)
}
