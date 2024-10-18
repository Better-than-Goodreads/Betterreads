package errors

import "net/http"

var (
	ErrPublishingBook = NewErrorDetails(
		"failed to publish book",
		"Error when publishing book: ",
		http.StatusBadRequest,
	)

	ErrParsingBookRequest = NewErrorDetails(
		"failed to parse request",
		"Error when parsing request: ",
		http.StatusBadRequest,
	)
)

func NewErrPublishingBook(err error) *ErrorDetails {
	errorDetails := NewErrorDetails(
		ErrPublishingBook.Title,
		ErrPublishingBook.Detail+err.Error(),
		ErrPublishingBook.Status,
	)
	return errorDetails

}

func NewErrParsingBookRequest(err error) *ErrorDetails {
	errorDetails := NewErrorDetails(
		ErrParsingBookRequest.Title,
		ErrParsingBookRequest.Detail+err.Error(),
		ErrParsingBookRequest.Status,
	)
	return errorDetails
}

func NewErrGettingBook(err error) *ErrorDetails {
	errorDetails := NewErrorDetails(
		"failed to get book",
		"Error when getting book: "+err.Error(),
		http.StatusBadRequest,
	)
	return errorDetails

}
