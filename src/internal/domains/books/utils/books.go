package utils

import (
	"github.com/betterreads/internal/domains/books/models"
	"github.com/google/uuid"
)

func MapBookRequestToBookRecord(req *models.NewBookRequest, id uuid.UUID, author uuid.UUID, authorName string) models.Book {
	return models.Book{
        Title:   req.Title,
        Author: author,
        AuthorName: authorName,
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


func MapBookDbToBook(book *models.BookDb, genres []string, ratings *models.Ratings, author string ) *models.Book{
    return &models.Book{
        Title: book.Title,
        Author: book.Author,
        AuthorName: author,
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



