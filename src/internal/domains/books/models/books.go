package models

import "github.com/google/uuid"

// RECORDS
type Book struct {
	Title           string    `json:"title" binding:"required" db:"title"`
	Author          uuid.UUID `json:"author" binding:"required" db:"author"`
	AuthorName      string    `json:"author_name" binding:"required" db:"author_name"`
	Description     string    `json:"description" binding:"required" db:"description"`
	AmountOfPages   int       `json:"amount_of_pages" binding:"required" db:"amount_of_pages"`
	PublicationDate string    `json:"publication_date" binding:"required" db:"publication_date"`
	Language        string    `json:"language" binding:"required" db:"language"`
	Genres          []string  `json:"genres" binding:"required" `
	Id              uuid.UUID `json:"id" binding:"required" db:"id"`
	TotalRatings    int       `json:"total_ratings" db:"total_ratings"`
	AverageRating   float64   `json:"avg_rating" db:"avg_rating"`
} // This struct is used to return the genres and book from the database

type BookRecord struct {
	Title           string    `json:"title" db:"title"`
	Author          uuid.UUID `json:"author" db:"author"`
	AuthorName      string    `json:"author" db:"author_name"`
	Description     string    `json:"description" db:"description"`
	AmountOfPages   int       `json:"amount_of_pages" db:"amount_of_pages"`
	PublicationDate string    `json:"publication_date" db:"publication_date"`
	Language        string    `json:"language" db:"language"`
	Id              uuid.UUID `json:"id" db:"id"`
	TotalRatings    int       `json:"total_ratings" db:"total_ratings"`
	AverageRating   float64   `json:"avg_rating" db:"avg_ratings"`
}

type BookDb struct {
	Title           string    `json:"title" db:"title"`
	Author          uuid.UUID `json:"author" db:"author"`
	Description     string    `json:"description" db:"description"`
	AmountOfPages   int       `json:"amount_of_pages" db:"amount_of_pages"`
	PublicationDate string    `json:"publication_date" db:"publication_date"`
	Language        string    `json:"language" db:"language"`
	Id              uuid.UUID `json:"id" db:"id"`
}

type GenreBook struct {
	GenreId int `json:"genre_id" db:"genre_id"`
	BookId  int `json:"book_id" db:"book_id"`
}

type Rating struct {
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	BookId uuid.UUID `json:"book_id" db:"book_id"`
	Rating int       `json:"rating" db:"rating"`
}

type Ratings struct {
	Total_ratings int     `json:"total_ratings" db:"total_ratings"`
	Avg_ratings   float64 `json:"avg_ratings" db:"avg_ratings"`
}

type ReviewDb struct {
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	BookId uuid.UUID `json:"book_id" db:"book_id"`
	Review string    `json:"review" db:"review"`
	Rating int       `json:"rating" db:"rating"`
}

type Review struct {
	Text   string `json:"review" db:"review"`
	Rating int    `json:"rating" db:"rating"`
}

// REQUESTS
type NewBookRequest struct {
	Title           string   `json:"title" validate:"required"`
	Description     string   `json:"description" validate:"required"`
	AmountOfPages   int      `json:"amount_of_pages" validate:"required"`
	PublicationDate string   `json:"publication_date" validate:"required"`
	Language        string   `json:"language" validate:"required"`
	Genres          []string `json:"genres" validate:"required"`
	Picture         []byte   `json:"picture"`
}

type NewRatingRequest struct {
	Rating int `json:"rating" binding:"required"`
}

type NewReviewRequest struct {
	Review string `json:"review"`
	Rating int    `json:"rating" binding:"required"`
}

// RESPONSES
type BookResponse struct {
	Title           string    `json:"title"`
	Author          uuid.UUID `json:"author_id"`
	AuthorName      string    `json:"author_name"`
	Description     string    `json:"description"`
	PublicationDate string    `json:"publication_date"`
	Language        string    `json:"language"`
	Genres          []string  `json:"genres"`
	AmountOfPages   int       `json:"amount_of_pages"`
	TotalRatings    int       `json:"total_ratings"`
	AverageRating   float64   `json:"avg_rating"`
	Id              uuid.UUID `json:"id"`
}

type ReviewOfUser struct {
	BookTitle string    `json:"book_title" db:"book_title"`
	Review    string    `json:"review" db:"review"`
	BookId    uuid.UUID `json:"book_id" db:"book_id"`
	Rating    int       `json:"rating" db:"rating"`
}

type ReviewOfBook struct {
	Username string    `json:"username" db:"username"`
	Review   string    `json:"review" db:"review"`
	UserId   uuid.UUID `json:"user_id" db:"user_id"`
	Rating   int       `json:"rating" db:"rating"`
}

type BookResponseWithReview struct {
	Book            *BookResponse `json:"book"`
	Review          *Review       `json:"review,omitempty"`
	BookShelfStatus *string       `json:"status,omitempty"`
}
