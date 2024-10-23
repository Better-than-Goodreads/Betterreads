package utils

import (
	"github.com/betterreads/internal/domains/books/models"
	"github.com/google/uuid"
)

func MapBookRequestToBookRecord(req *models.NewBookRequest, id uuid.UUID, author uuid.UUID) models.Book {
	return models.Book{
        Title:   req.Title,
        Author: author,
        Description: req.Description,
        AmountOfPages: req.AmountOfPages,
        PublicationDate: req.PublicationDate,
        Language: req.Language,
        Genres: req.Genres,
        Id: id,
    }
}


func MapBookToBookResponse(book *models.Book, author_name string) *models.BookResponse {
    return &models.BookResponse{
        Title: book.Title,
        Author: author_name, 
        Description: book.Description,
        PublicationDate: book.PublicationDate,
        Language: book.Language,
        Genres: book.Genres,
        AmountOfPages: book.AmountOfPages,
        Id: book.Id,
        TotalRatings: book.TotalRatings,
        AverageRating: book.AverageRating,
    }
}


func MapBookDbToBook(book *models.BookDb, genres []string, ratings *models.Ratings ) *models.Book{
    return &models.Book{
        Title: book.Title,
        Author: book.Author,
        Description: book.Description,
        PublicationDate: book.PublicationDate,
        Language: book.Language,
        Genres: genres,
        AmountOfPages: book.AmountOfPages,
        Id: book.Id,
        TotalRatings: ratings.Total_ratings,
        AverageRating: ratings.Avg_ratings,
    }
}


func MapRatingToRatingResponse (rating *models.Rating) *models.RatingResponse{
    return &models.RatingResponse{
        BookId: rating.BookId,
        Rating: rating.Rating,
    }
}
