package models

import "github.com/google/uuid"

type BookInShelf struct {
	Date   string    `json:"date" db:"date"`
	Status string    `json:"status" db:"status"`
	BookId uuid.UUID `json:"book_id" db:"book_id"`
	UserId uuid.UUID `json:"user_id" db:"user_id"`
}

type BookShelfBasicData struct {
	Date       string    `json:"date" db:"date"`
	Status     string    `json:"status" db:"status"`
	Title      string    `json:"title" db:"title"`
	AuthorId   string    `json:"author_id" db:"author_id"`
	AuthorName string    `json:"author_name" db:"author_name"`
	BookId     uuid.UUID `json:"book_id" db:"book_id"`
	UserId     uuid.UUID `json:"user_id" db:"user_id"`
}

type BookShelfUserReview struct {
	BookId uuid.UUID `json:"book_id" db:"book_id"`
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	Review string    `json:"review" db:"review"`
	Rating int       `json:"rating" db:"rating"`
}

type BookShelfRatingStats struct {
	BookId       uuid.UUID `json:"book_id" db:"book_id"`
	UserId       uuid.UUID `json:"user_id" db:"user_id"`
	AvgRating    float64   `json:"avg_ratings" db:"avg_ratings"`
	TotalRatings int       `json:"total_ratings" db:"total_ratings"`
}

type BookInShelfResponse struct {
	Date         string    `json:"date_added" db:"date"`
	Status       string    `json:"status" db:"status"`
	Title        string    `json:"title" db:"title"`
	AuthorId     string    `json:"author" db:"author_id"`
	AuthorName   string    `json:"author_name" db:"author_name"`
	UserReview   string    `json:"user_review" db:"user_review"`
	AvgRating    float64   `json:"avg_ratings" db:"avg_ratings"`
	TotalRatings int       `json:"total_ratings" db:"total_ratings"`
	BookId       uuid.UUID `json:"book_id" db:"book_id"`
	UserRating   int       `json:"user_rating" db:"user_rating"`
}

type BookShelfRequest struct {
	Status string    `json:"status" binding:"required"`
	BookId uuid.UUID `json:"book_id" binding:"required"`
}

type BookShelfType string

const (
	BookShelfTypeWantToRead BookShelfType = "plan-to-read"
	BookShelfTypeReading    BookShelfType = "reading"
	BookShelfTypeRead       BookShelfType = "read"
	BookShelfAll            BookShelfType = "all"
)

var ValidBookShelfTypes = []BookShelfType{BookShelfTypeWantToRead, BookShelfTypeReading, BookShelfTypeRead, BookShelfAll}
