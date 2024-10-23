package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	gin "github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	ErrParsingRequest = NewErrorDetails(
		"failed to parse request",
		"Error when parsing request: ",
		http.StatusBadRequest,
	)
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

type ErrorDetailsWithParams struct {
	Type     string       `json:"type"`
	Title    string       `json:"title"`
	Detail   string       `json:"detail"`
	Instance string       `json:"instance"`
	Params   []ErrorParam `json:"validation_errors"`
	Status   int          `json:"status"`
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
			fmt.Println(err)
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

func NewErrorDetailsWithParams(title string, detail string, status int, err error) *ErrorDetailsWithParams {
	return &ErrorDetailsWithParams{
		Type:   "about:blank",
		Title:  title,
		Detail: detail,
		Status: status,
		Params: parseParameters(err),
	}
}

func NewErrParsingRequest(err error) *ErrorDetailsWithParams {
	errorDetails := NewErrorDetailsWithParams(
		ErrParsingRequest.Title,
		ErrParsingRequest.Detail,
		ErrParsingRequest.Status,
		err,
	)
	return errorDetails
}

func NewErrNotLogged() *ErrorDetails {
	errNotLogged := NewErrorDetails(
		"not logged",
		"User is not logged",
		http.StatusUnauthorized,
	)
	return errNotLogged
}

func SendError(c *gin.Context, err *ErrorDetails) {
	err.Instance = c.Request.RequestURI
	c.JSON(err.Status, err)
}

func SendErrorWithParams(c *gin.Context, err *ErrorDetailsWithParams) {
	err.Instance = c.Request.RequestURI
	c.JSON(err.Status, err)
}
