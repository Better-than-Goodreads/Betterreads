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

func NewErrGettingBookReviews(err error) *ErrorDetails {
    errGettingReviews := NewErrorDetails(
        "failed to get reviews",
        "Error when getting reviews: " + err.Error(),
        http.StatusInternalServerError,
    )
    return errGettingReviews
}

func NewErrBookNotFound() *ErrorDetails {
	errBookNotFound := NewErrorDetails(
		"book not found",
		"Book not found",
		http.StatusNotFound,
	)
	return errBookNotFound
}

func NewErrAuthorNotFound() *ErrorDetails {
    errAuthorNotFound := NewErrorDetails(
        "author not found",
        "Author not found",
        http.StatusNotFound,
    )
    return errAuthorNotFound
}

func NewErrNotAuthor() *ErrorDetails {
    errNotAuthor := NewErrorDetails(
        "not an author",
        "User is not an author",
        http.StatusUnauthorized,
    )
    return errNotAuthor
}

func NewErrInvalidRatingBook(err error) *ErrorDetailsWithParams {
    errRatingBook := NewErrorDetailsWithParams(
		"failed to rate book",
		"Rate value invalid",
		http.StatusBadRequest,
        err,
	)
	return errRatingBook
}

func NewErrGettingBooks(err error) *ErrorDetails {
    errGettingBooks := NewErrorDetails(
        "failed to get books",
        "Error when getting books: " + err.Error(),
        http.StatusInternalServerError,
    )
    return errGettingBooks
}

func NewErrInvalidBookId(id string) *ErrorDetails {
    errInvalidBookId := NewErrorDetails(
		"failed to parse book id",
		"value of id should be a uuid: " + "Id: " + id,
		http.StatusBadRequest,
	)

    return errInvalidBookId
}

func NewErrInvalidAuthorId(id string) *ErrorDetails {
    errInvalidAuthorId := NewErrorDetails(
        "failed to parse author id",
        "value of id should be a uuid: " + "Id: " + id,
        http.StatusBadRequest,
    )
    return errInvalidAuthorId
}

func NewErrInvalidRating(rate string) *ErrorDetails {
    errorDetails := NewErrorDetails(
        "failed to parse rate",
        "Rate should be a number: "+rate,
        http.StatusBadRequest,
    )
    return errorDetails
}

func NewErrRatingNotFound() *ErrorDetails {
	errRatingNotFound := NewErrorDetails(
		"rating not found",
		"Rating not found",
		http.StatusNotFound,
	)
	return errRatingNotFound
}

func NewErrGettingRating(err error) *ErrorDetails {
	errGettingRating := NewErrorDetails(
		"failed to get rating",
		"Error when getting rating: " + err.Error(),
		http.StatusInternalServerError,
	)
	return errGettingRating
}

func NewErrDeletingRating(err error) *ErrorDetails {
	errDeletingRating := NewErrorDetails(
		"failed to delete rating",
		"Error when deleting rating: " + err.Error(),
		http.StatusInternalServerError,
	)
	return errDeletingRating
}

func NewErrAddingReview(err error) *ErrorDetails {
	errAddingReview := NewErrorDetails(
		"failed to add review",
		"Error when adding review: " + err.Error(),
		http.StatusInternalServerError,
	)
	return errAddingReview
}


func NewErrBookNotFoundByName(name string) *ErrorDetails {
    err := NewErrorDetails(
        "failed to get book",
        "Book not found with name " + name,
        http.StatusNotFound,
    )
    return err
}

func NewErrNoPicture() *ErrorDetailsWithParams{
    err := ErrorParam{
        Name: "picture",
        Reason: "picture not found",
    }
    errNoPicture := NewErrorDetailsWithParams(
        "failed to get book picture",
        "Book picture not found",
        http.StatusNotFound,
        err,
    )
    return errNoPicture
}

func NewErrRatingAlreadyExists() *ErrorDetails {
    err := NewErrorDetails(
        "failed to rate book",
        "Rating already exists",
        http.StatusBadRequest,
    )
    return err
}


func NewErrRating(err error) *ErrorDetailsWithParams {
    errUpdatingRating := NewErrorDetailsWithParams(
        "failed to rate",
        "Error when rating: " + err.Error(),
        http.StatusInternalServerError,
        err,
    )
    return errUpdatingRating
}

func NewErrReviewAlreadyExists() *ErrorDetails {
    err := NewErrorDetails(
        "failed to add review",
        "Review already exists",
        http.StatusBadRequest,
    )
    return err
}

