package errors

import (
	"encoding/json"

	gin "github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// Follows RFC 7807: https://datatracker.ietf.org/doc/html/rfc7807
type ErrorDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
	Status   int    `json:"status"` // ESTE NO IRIA MAS
}

func (e ErrorDetails) Error() string {
	string := e.Title + " " + e.Detail
	return string
}

func NewErrorDetails(title string, err error, status int) *ErrorDetails {
	return &ErrorDetails{
		Type:   "about:blank",
		Title:  title,
		Detail: err.Error(),
		Status: status,
	}
}

type ErrorDetailsWithParams struct {
	Type     string       `json:"type"`
	Title    string       `json:"title"`
	Detail   string       `json:"detail"`
	Instance string       `json:"instance"`
	Params   []ErrorParam `json:"validation_errors"`
	Status   int          `json:"status"`
}

func (e ErrorDetailsWithParams) Error() string {
	string := e.Title + " " + e.Detail
	return string
}

type ErrorParam struct {
	Name   string `json:"field"`
	Reason string `json:"reason"`
}

func (e ErrorParam) Error() string {
	return e.Reason
}

func parseParameters(err error) []ErrorParam {
	var errors []ErrorParam

	if errParam, ok := err.(ErrorParam); ok {
		errors = append(errors, errParam)
	} else if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
		errors = append(errors, ErrorParam{
			Name:   unmarshalErr.Field,
			Reason: unmarshalErr.Type.String(),
		})
	} else if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			errors = append(errors, ErrorParam{
				Name:   err.Field(),
				Reason: err.Tag(),
			})
		}
	} else {
		errors = append(errors, ErrorParam{
			Name:   "unknown",
			Reason: err.Error(),
		})
	}
	return errors
}

func NewErrorDetailsWithParams(title string, status int, err error) *ErrorDetailsWithParams {
	return &ErrorDetailsWithParams{
		Type:   "about:blank",
		Title:  title,
		Detail: "Invalid request parameters",
		Status: status,
		Params: parseParameters(err),
	}
}

func AbortWithJsonErorr(c *gin.Context, err error) {
	errToSend := NewErrorDetailsWithParams("Error parsing request json", http.StatusBadRequest, err)
	c.AbortWithError(errToSend.Status, errToSend)
}

func SendError(c *gin.Context, err *ErrorDetails) {
	err.Instance = c.Request.RequestURI
	c.JSON(err.Status, err)
}

func SendErrorWithParams(c *gin.Context, err *ErrorDetailsWithParams) {
	err.Instance = c.Request.RequestURI
	c.JSON(err.Status, err)
}
