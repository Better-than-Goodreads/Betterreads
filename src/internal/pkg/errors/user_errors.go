package errors

import (
	"net/http"
)

var (
	ErrFetchUsers = NewErrorDetails(
		"failed to fetch all users",
		"Error when Fetching all users:  ",
		http.StatusInternalServerError,
	)

	ErrRegisterUser = NewErrorDetails(
		"failed to register user",
		"Error when registering user: ",
		http.StatusBadRequest,
	)

	ErrLogInUser = NewErrorDetails(
		"failed to log in user",
		"Error when logging in user: ",
		http.StatusBadRequest,
	)


	ErrInvalidID = NewErrorDetails(
		"failed to parse user id",
		"value of id should be a number: ",
		http.StatusBadRequest,
	)

	ErrUserNotFoundById = NewErrorDetails(
		"failed to fetch user by id",
		"Error when fetching user by id: ",
		http.StatusNotFound,
	)
)

func NewErrFetchUsers(err error) *ErrorDetails {
	errorDetails := NewErrorDetails(
		ErrFetchUsers.Title,
		ErrFetchUsers.Detail+err.Error(),
		ErrFetchUsers.Status,
	)
	return errorDetails
}

func NewErrRegisterUser(err error) *ErrorDetails {
	errorDetails := NewErrorDetails(
		ErrRegisterUser.Title,
		ErrRegisterUser.Detail+err.Error(),
		ErrRegisterUser.Status,
	)
	return errorDetails
}

func NewErrLogInUser(err error) *ErrorDetails {
	errorDetails := NewErrorDetails(
		ErrLogInUser.Title,
		ErrLogInUser.Detail+err.Error(),
		ErrLogInUser.Status,
	)
	return errorDetails
}

func NewErrInvalidID(id string) *ErrorDetails {
	errorDetails := NewErrorDetails(
		ErrInvalidID.Title,
		ErrInvalidID.Detail+"Id: "+id,
		ErrInvalidID.Status,
	)
	return errorDetails
}

func NewErrUserNotFoundById(err error) *ErrorDetails {
	errorDetails := NewErrorDetails(
		ErrUserNotFoundById.Title,
		ErrUserNotFoundById.Detail+err.Error(),
		ErrUserNotFoundById.Status,
	)
	return errorDetails
}
