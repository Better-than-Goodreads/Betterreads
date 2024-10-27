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
		http.StatusInternalServerError,
	)
	
	return errRegisterUser 
}

func NewErrUserNotUnique(err error) *ErrorDetailsWithParams{
	errUserNotUnique := NewErrorDetailsWithParams(
		"failed to register user",
		"Error when registering user: " + err.Error(),
		http.StatusBadRequest,
		err,
	)

	return errUserNotUnique
}

func NewErrLogInUser(err error) *ErrorDetails {
    errLogInUser := NewErrorDetails(
		"failed to log in user",
		"Error when logging in user: " + err.Error(),
		http.StatusUnauthorized,
	)

	return errLogInUser
}

func NewErrInvalidRegisterId(id string) *ErrorDetails {
    errInvalidRegisterId := NewErrorDetails(
        "Invalid id",
        "Id: " + id + " is not a valid uuid",
        http.StatusBadRequest,
    )
    return errInvalidRegisterId
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

func NewErrParsingPicture() *ErrorDetailsWithParams {
    err := ErrorParam {
        Name: "file",
        Reason: "file is not a valid picture / file is empty",
    }
    errParsingPicture := NewErrorDetailsWithParams(
        "failed to parse picture",
        "Error when parsing picture",
        http.StatusBadRequest,
        err,
    )
    return errParsingPicture
}


func NewErrNoPictureUser() *ErrorDetails{
    errNoPicture := NewErrorDetails(
        "failed to get user picture",
        "user picture not found",
        http.StatusNotFound,
    )
    return errNoPicture
}


func NewErrPostPicture() *ErrorDetails{
    errPostPicture := NewErrorDetails(
        "failed to post user picture",
        "Error when posting user picture",
        http.StatusInternalServerError,
    )
    return errPostPicture
}
