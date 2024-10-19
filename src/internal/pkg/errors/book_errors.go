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

	ErrRatingBook = NewErrorDetails(
		"failed to rate book",
		"Rate must be between 1 and 5 (inclusive): ",
		http.StatusBadRequest,
	)

	ErrorGettingBook = NewErrorDetails(
		"failed to get book",
		"Error when getting book: ",
		http.StatusBadRequest,
	)

	ErrorParsingBookRequest = NewErrorDetails(
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

func NewErrBookNotFound() *ErrorDetails {
	errorDetails := NewErrorDetails(
		"book not found",
		"Book not found",
		http.StatusNotFound,
	)
	return errorDetails
}

func NewErrRatingBook(err error) *ErrorDetails {
	errorDetails := NewErrorDetails(
		ErrRatingBook.Title,
		ErrRatingBook.Detail+err.Error(),
		ErrRatingBook.Status,
	)
	return errorDetails
}

func NewErrParsingError(err error) *ErrorDetails {
	errorDetails := NewErrorDetails(
		"failed to parse error",
		"Error when parsing error: "+err.Error(),
		http.StatusBadRequest,
	)
	return errorDetails
}