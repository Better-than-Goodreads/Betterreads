package errors

import "net/http"

func NewErrPublishingBook(err error) *ErrorDetails {
    ErrPublishingBook := NewErrorDetails(
		"failed to publish book",
		"Error when publishing book: " + err.Error(),
		http.StatusBadRequest,
	)
	return ErrPublishingBook

}

func NewErrParsingBookRequest(err error) *ErrorDetails {
    errParsingBookRequest := NewErrorDetails(
		"failed to parse request",
		"Error when parsing request: " + err.Error(),
		http.StatusBadRequest,
	)
    return errParsingBookRequest
}

func NewErrGettingBook(err error) *ErrorDetails {
    errGettingBooks := NewErrorDetails(
		"failed to get book",
		"Error when getting book: " + err.Error(),
		http.StatusBadRequest,
	)
    return errGettingBooks
}

func NewErrBookNotFound() *ErrorDetails {
	errBookNotFound := NewErrorDetails(
		"book not found",
		"Book not found",
		http.StatusNotFound,
	)
	return errBookNotFound
}

func NewErrRatingBook(err error) *ErrorDetails {
    errRatingBook := NewErrorDetails(
		"failed to rate book",
		"Rate must be between 1 and 5 (inclusive): "+ err.Error(),
		http.StatusBadRequest,
	)

	return errRatingBook
}

func NewErrInvalidBookId(id string) *ErrorDetails {
    errInvalidBookId := NewErrorDetails(
		"failed to parse book id",
		"value of id should be a number: " + "Id: " + id,
		http.StatusBadRequest,
	)

    return errInvalidBookId
}

func NewErrInvalidRating(rate string) *ErrorDetails {
    errorDetails := NewErrorDetails(
        "failed to parse rate",
        "Rate should be a number: "+rate,
        http.StatusBadRequest,
    )
    return errorDetails
}

