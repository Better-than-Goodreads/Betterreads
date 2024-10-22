package utils

import (
	"github.com/betterreads/internal/domains/books/models"
	"github.com/google/uuid"
)

func MapBookRequestToBookRecord(req *models.NewBookRequest, id uuid.UUID, author string) models.Book {
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


func MapBookToBookResponse(book *models.Book) *models.BookResponse {
    return &models.BookResponse{
        Title: book.Title,
        Author: book.Author,
        Description: book.Description,
        PublicationDate: book.PublicationDate,
        Language: book.Language,
        Genres: book.Genres,
        AmountOfPages: book.AmountOfPages,
        Id: book.Id,
    }
}


func MapBookDbToBook(book *models.BookDb, genres []string) *models.Book{
    return &models.Book{
        Title: book.Title,
        Author: book.Author,
        Description: book.Description,
        PublicationDate: book.PublicationDate,
        Language: book.Language,
        Genres: genres,
        AmountOfPages: book.AmountOfPages,
        Id: book.Id,
    }
}
