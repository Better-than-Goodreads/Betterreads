package model

import (
	"github.com/betterreads/internal/domains/books/models"
)

type RecommendationsByGenre struct{
    Genre string `json:"genre"`
    Books []*models.Book `json:"books"`
}
