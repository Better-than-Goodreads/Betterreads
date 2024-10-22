package repository

import (
    "github.com/google/uuid"
)

type Book struct {
	Title           string   `json:"title" binding:"required"`
	Author          string   `json:"author" binding:"required"`
	Description     string   `json:"description" binding:"required"`
	PhotoId         string   `json:"photo_id" binding:"required"`
	AmountOfPages   string   `json:"amount_of_pages" binding:"required"`
	PublicationDate string   `json:"publication_date" binding:"required"`
	Language        string   `json:"language" binding:"required"`
	Genres          []string `json:"genres" binding:"required"`
	Ratings         map[int]int    `json:"ratings"` //Tal vez haya que modificar esto mas adelante
	// El id de un rating es IdBookIdUser, los 2 numeros concatenados
}

type BookDb struct {
	Id 				uuid.UUID      `json:"id" db:"id"`
	Title           string   `json:"title" db:"title"`
	Author          string   `json:"author" db:"author"`
	Description     string   `json:"description" db:"description"`
	AmountOfPages   string   `json:"amount_of_pages" db:"amount_of_pages"`
	
	// Ratings         map[int]int    `json:"ratings"` //Tal vez haya que modificar esto mas adelante
	// El id de un rating es IdBookIdUser, los 2 numeros concatenados
}

type BooksDatabase interface {
	SaveBook(book Book) error
	GetBookById(id int) (*Book, error)
	GetBookByName(name string) (*Book, error)
	RateBook(bookId int, userId int, rating int) error
	DeleteRating(bookId int, userId int) error
}
