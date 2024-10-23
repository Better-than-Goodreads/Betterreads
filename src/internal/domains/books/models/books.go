package models

import "github.com/google/uuid"

// RECORDS
type Book struct {
	Title           string    `json:"title" binding:"required"`
	Author          uuid.UUID    `json:"author" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	AmountOfPages   int       `json:"amount_of_pages" binding:"required"`
	PublicationDate string    `json:"publication_date" binding:"required"`
	Language        string    `json:"language" binding:"required"`
	Genres          []string  `json:"genres" binding:"required"`
	Id              uuid.UUID `json:"id" binding:"required"`
	TotalRatings    int       `json:"total_ratings"`
	AverageRating   float64   `json:"avg_rating"`
} // This struct is used to return the genres and book from the database

type BookDb struct {
	Title           string    `json:"title" db:"title"`
	Author          uuid.UUID    `json:"author" db:"author"`
	Description     string    `json:"description" db:"description"`
	AmountOfPages   int       `json:"amount_of_pages" db:"amount_of_pages"`
	PublicationDate string    `json:"publication_date" db:"publication_date"`
	Language        string    `json:"language" db:"language"`
	Id              uuid.UUID `json:"id" db:"id"`
	// Ratings         map[int]int    `json:"ratings"` //Tal vez haya que modificar esto mas adelante
	// El id de un rating es IdBookIdUser, los 2 numeros concatenados
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

// REQUESTS
type NewBookRequest struct {
	Title           string   `json:"title" validate:"required"`
	Description     string   `json:"description" validate:"required"`
	AmountOfPages   int      `json:"amount_of_pages" validate:"required"`
	PublicationDate string   `json:"publication_date" validate:"required"`
	Language        string   `json:"language" validate:"required"`
	Genres          []string `json:"genres" validate:"required"`
	Picture		 	[]byte	 `json:"picture"`
}

type NewRatingRequest struct {
	Rating int `json:"rating" binding:"required"`
}

// RESPONSES
type BookResponse struct {
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	Description     string    `json:"description"`
	PublicationDate string    `json:"publication_date"`
	Language        string    `json:"language"`
	Genres          []string  `json:"genres"`
	AmountOfPages   int       `json:"amount_of_pages"`
	TotalRatings    int       `json:"total_ratings"`
	AverageRating   float64   `json:"avg_rating"`
	Id              uuid.UUID `json:"id"`
}

type RatingResponse struct {
	BookId uuid.UUID `json:"book_id"`
	Rating int       `json:"rating"`
}
