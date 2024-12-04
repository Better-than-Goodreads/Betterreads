package models

import (
	"github.com/google/uuid"
)

type Post struct {
	UserId          uuid.UUID `json:"id" db:"user_id"`
	Username        string    `json:"username" db:"username"`
	BookId          uuid.UUID `json:"book_id" db:"book_id"`
	BookAuthor      string    `json:"book_author" db:"author_name"`
	BookTitle       string    `json:"book_title" db:"title"`
	BookDescription string    `json:"book_description" db:"description"`
	PublicationDate string    `json:"publication_date" db:"publication_date"`
	// Ratings can be null to have the two posts.
	// - The publication of a book
	// - The rating of a book
	Rating *int `json:"rating,omitempty" db:"rating"`
}

type PostDTO struct {
	Type string `json:"type"`
	Post Post   `json:"post"`
}
