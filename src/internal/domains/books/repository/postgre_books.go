package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"fmt"
	"strings"
)

type PostgresBookRepository struct {
	c *sqlx.DB
}

func NewPostgresBookRepository(c *sqlx.DB) (BooksDatabase, error) {
	enableUUIDExtension := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	if _, err := c.Exec(enableUUIDExtension); err != nil {
		return nil, fmt.Errorf("failed to enable uuid extension: %w", err)
	}

	schemaBooks := `
		CREATE TABLE IF NOT EXISTS books (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL,
			description VARCHAR(255) NOT NULL,
			photo_id VARCHAR(255),
			amount_of_pages INTEGER NOT NULL,
			publication_date VARCHAR(255) NOT NULL,
			language VARCHAR(255),
			genres VARCHAR(255)
			);			
			`			
			// CREATE UNIQUE INDEX IF NOT EXISTS idx_books_title_author ON books(title, author);
			// ratings_table UUID
			// FOREIGN KEY (ratings_table) REFERENCES ratings(id)
	

	schemaRatings := `
		CREATE TABLE IF NOT EXISTS ratings (
			user_id UUID,
			book_id UUID,
			rating INTEGER,
			PRIMARY KEY (user_id, book_id),
			FOREIGN KEY (book_id) REFERENCES books(id)			
			);
			`
			// FOREIGN KEY (user_id) REFERENCES users(id),

	// schemaBooks := `
	// 	DROP TABLE books
	// `

	if _, err := c.Exec(schemaBooks); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	if _, err := c.Exec(schemaRatings); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PostgresBookRepository{c}, nil
}

func (r *PostgresBookRepository) SaveBook(book Book) error {
	query := `INSERT INTO books (title, author, description, amount_of_pages, publication_date, language, genres)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	args := []interface{}{book.Title, book.Author, book.Description, book.AmountOfPages, book.PublicationDate, book.Language, strings.Join(book.Genres,",")}

	_, err := r.c.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	return nil
}
func (r *PostgresBookRepository) GetBookById(id int) (*Book, error) {
	return nil, nil
}
func (r *PostgresBookRepository) GetBookByName(name string) (*Book, error) {
	bookdb := &BookDb{}

	// Ratings         map[int]int    `json:"ratings"` //Tal vez haya que modificar esto mas adelante
	// El id de un rating es IdBookIdUser, los 2 numeros concatenados

	query:= `SELECT id, title, author, description, amount_of_pages FROM books 
			WHERE title = $1;`
	if err := r.c.Get(bookdb, query, name); err != nil {
		if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user with id %s not found", name)
        }
		return nil, fmt.Errorf("failed to get the book: %w", err)
	}

	book := &Book{
		
		Title: bookdb.Title,
		Author: bookdb.Author,
		Description: bookdb.Description,
		AmountOfPages: bookdb.Id.String(),

	}

	return book, nil
}
func (r *PostgresBookRepository) RateBook(bookId int, userId int, rating int) error {
	query := `INSERT INTO ratings (user_id, book_id, rating)
	VALUES ($1, $2, $3)
	`
	args := []interface{}{userId, bookId, rating}

	_, err := r.c.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to rate book: %w", err)
	}
	return nil
}
func (r *PostgresBookRepository) DeleteRating(bookId int, userId int) error {
	
	query := `DELETE FROM ratings WHERE user_id = $1 AND book_id = $2;`
	
	if _, err := r.c.Exec(query, userId, bookId); err != nil {
		return fmt.Errorf("failed to delete rating: %w", err)
	}
	
	return nil
}

func (r *PostgresBookRepository) GetRatings(bookId int, userId int) (*Rating, error) {
	var ratings *Rating
	query := `SELECT * FROM ratings WHERE book_id = $1 AND user_id = $2;`
	args := []interface{}{bookId, userId}

	if err := r.c.Select(&ratings, query, args); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no ratings found")
		}
		return nil, fmt.Errorf("failed to get ratings: %w", err)
	}
	return ratings, nil
}

