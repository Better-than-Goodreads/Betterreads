package utils

import (
	"github.com/betterreads/internal/domains/books/models"
	"github.com/betterreads/internal/domains/books/repository"
)

func MapBookRequestToBookRecord(req *models.NewBookRequest) repository.Book {
	return repository.Book{
		Title:           req.Title,
		Author:          req.Author,
		Description:     req.Description,
		PhotoId:         req.PhotoId,
		AmountOfPages:   req.AmountOfPages,
		PublicationDate: req.PublicationDate,
		Language:        req.Language,
		Genres:          req.Genres,
	}
}
