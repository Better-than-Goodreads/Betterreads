package model

import (
	"github.com/betterreads/internal/domains/books/models"
	"github.com/google/uuid"
)

type RecommendationsByGenre struct {
	Genre string         `json:"genre"`
	Books []*models.Book `json:"books"`
}

type BookRecommendation struct {
	Title           string    `json:"title" binding:"required" db:"title"`
	Author          uuid.UUID `json:"author" binding:"required" db:"author"`
	AuthorName      string    `json:"author_name" binding:"required" db:"author_name"`
	Description     string    `json:"description" binding:"required" db:"description"`
	AmountOfPages   int       `json:"amount_of_pages" binding:"required" db:"amount_of_pages"`
	PublicationDate string    `json:"publication_date" binding:"required" db:"publication_date"`
	Language        string    `json:"language" binding:"required" db:"language"`
	Id              uuid.UUID `json:"id" binding:"required" db:"id"`
	TotalRatings    int       `json:"total_ratings" db:"total_ratings"`
	AverageRating   float64   `json:"avg_rating" db:"avg_ratings"`
} // This struct is used to get all the values from the db to then sort them by rating. The genres are not in here because we are dummies that still use a hashmap.
