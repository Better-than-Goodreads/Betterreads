package errors

import (
	"net/http"
)


func NewErrFetchUsers(err error) *ErrorDetails {
    errFetchUsers := NewErrorDetails(
		"failed to fetch all users",
		"Error when Fetching all users:  " + err.Error(),
		http.StatusInternalServerError,
	)

	return errFetchUsers
}

func NewErrRegisterUser(err error) *ErrorDetails {
    errRegisterUser := NewErrorDetails(
		"failed to register user",
		"Error when registering user: " + err.Error(),
		http.StatusBadRequest,
	)
	
	return errRegisterUser 
}

func NewErrLogInUser(err error) *ErrorDetails {
    errLogInUser := NewErrorDetails(
		"failed to log in user",
		"Error when logging in user: " + err.Error(),
		http.StatusBadRequest,
	)

	return errLogInUser
}

func NewErrInvalidUserID(id string) *ErrorDetails {
    errInvalidUserID := NewErrorDetails(
		"failed to parse user id",
        "value of id should be a number: Id: " + id,
		http.StatusBadRequest,
	)

	return errInvalidUserID
}

func NewErrUserNotFoundById(err error) *ErrorDetails {
    errUserNotFoundById := NewErrorDetails(
		"failed to fetch user by id",
		"Error when fetching user by id: " + err.Error(),
		http.StatusNotFound,
	)

	return errUserNotFoundById
}
