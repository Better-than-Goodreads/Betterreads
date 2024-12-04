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
	Title           string    `json:"title" db:"title"`
	Author          string    `json:"author" db:"author_id"`
	AuthorName      string    `json:"author_name" db:"author_name"`
	Description     string    `json:"description" db:"description"`
	PublicationDate string    `json:"publication_date" db:"publication_date"`
	Date            string    `json:"date_added" db:"date"`
	Language        string    `json:"language" db:"language"`
	Genres          string    `json:"ignore,omitempty" db:"genres"`
	GenresArray     *[]string `json:"genres" db:"genres_array"`
	AmountOfPages   int       `json:"amount_of_pages" db:"amount_of_pages"`
	TotalRatings    int       `json:"total_ratings" db:"total_ratings"`
	AvgRating       float64   `json:"avg_ratings" db:"avg_ratings"`
	Status          string    `json:"status" db:"status"`
	UserReview      string    `json:"user_review" db:"user_review"`
	UserRating      int       `json:"user_rating" db:"user_rating"`
	Id              uuid.UUID `json:"book_id" db:"id"`
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
