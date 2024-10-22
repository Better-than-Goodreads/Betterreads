package repository

import (
	"database/sql"

	"fmt"

	"github.com/betterreads/internal/domains/books/models"
	 "github.com/betterreads/internal/domains/books/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresBookRepository struct {
	c *sqlx.DB
}

func getGenreById(genre string) (int, error) {
    for key, value := range genresDict {
        if value == genre {
            return key, nil
        }
    }
    return -1 , ErrGenreNotFound
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
			amount_of_pages INTEGER NOT NULL,
			publication_date VARCHAR(255) NOT NULL,
            language VARCHAR(255) NOT NULL
			);
			`
			// CREATE UNIQUE INDEX IF NOT EXISTS idx_books_title_author ON books(title, author);
			// ratings_table UUID
			// FOREIGN KEY (ratings_table) REFERENCES ratings(id)
    

    schemaGendersBooks := `
        CREATE TABLE IF NOT EXISTS genres_books (
            book_id UUID,
            genre_id INT,
            PRIMARY KEY (book_id, genre_id),
            FOREIGN KEY (book_id) REFERENCES books(id)
        );
    `
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

    if _, err := c.Exec(schemaGendersBooks); err != nil {
        return nil, fmt.Errorf("failed to create table: %w", err)
    }
	return &PostgresBookRepository{c}, nil
}

func (r *PostgresBookRepository) SaveBook (book *models.NewBookRequest, author string) (*models.Book, error) {
    bookRecord := &models.BookDb{}
    query := `INSERT INTO books (title, author, description,  amount_of_pages,
                    publication_date, language)
                    VALUES ($1, $2, $3, $4, $5, $6)
                    RETURNING id, title, author, description, amount_of_pages, publication_date, language;`

    args := []interface{}{book.Title, author, book.Description, book.AmountOfPages, book.PublicationDate, book.Language}

    if err := r.c.Get(bookRecord, query, args...); err != nil {
        return nil , fmt.Errorf("failed to create book: %w", err)
    }
    
    query = `INSERT INTO genres_books (book_id, genre_id)
             VALUES ($1, $2);`

    for _, genre := range book.Genres {
        genreid, err := getGenreById(genre)
        if err != nil {
            return nil, fmt.Errorf("failed to create book: %w", err)
        }
        args = []interface{}{bookRecord.Id, genreid}
        if _, err := r.c.Exec(query, args...); err != nil {
            return nil, fmt.Errorf("failed to create book: %w", err)
        }
    }

    res := utils.MapBookRequestToBookRecord(book, bookRecord.Id, author)

    return &res, nil
}

func (r *PostgresBookRepository) getGenresForBook(book_id uuid.UUID) ([]string, error) {
    var genres_ids []int
    query := `SELECT genre_id FROM genres_books WHERE book_id = $1;`
    if err := r.c.Select(&genres_ids, query, book_id); err != nil {
        return nil, fmt.Errorf("failed to get genres: %w", err)
    }

    genres := []string{}
    for _, genre_id := range genres_ids {
        genres = append(genres, genresDict[genre_id])
    }
    
    return genres, nil

}

func (r *PostgresBookRepository) GetBookById(id uuid.UUID) (*models.Book, error) {
    bookdb := &models.BookDb{}
    query := `SELECT * FROM books WHERE id = $1;`
    if err := r.c.Get(bookdb, query, id); err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("book with id %s not found", id)
        }
        return nil, fmt.Errorf("failed to get book: %w", err)
    }
    genres, err := r.getGenresForBook(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get book: %w", err)
    }
    
    book := utils.MapBookDbToBook(bookdb, genres)

    return book, nil
}

func (r *PostgresBookRepository) GetBooks() ([]*models.Book, error) {
    var books []*models.BookDb
    query := `SELECT * FROM books;`
    if err := r.c.Select(&books, query); err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("No books found")
        }
        return nil, fmt.Errorf("failed to get books: %w", err)
    }
    res := []*models.Book{}
    for _, book := range books {
        genres, err := r.getGenresForBook(book.Id)
        if err != nil {
            return nil, fmt.Errorf("failed to get books: %w", err)
        }
        res = append(res, utils.MapBookDbToBook(book, genres))
    }
    return res, nil
} 


func (r *PostgresBookRepository) GetBookByName(name string) (*models.Book, error) {
	bookdb := &models.BookDb{}

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

	book := &models.Book{
		//
		// Title: bookdb.Title,
		// Author: bookdb.Author,
		// Description: bookdb.Description,
		// AmountOfPages: bookdb.Id.String(),
		//
	}

	return book, nil
}
func (r *PostgresBookRepository) RateBook(bookId uuid.UUID, userId uuid.UUID, rating int) error {
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
func (r *PostgresBookRepository) DeleteRating(bookId uuid.UUID, userId uuid.UUID) error {
	
	query := `DELETE FROM ratings WHERE user_id = $1 AND book_id = $2;`
	
	if _, err := r.c.Exec(query, userId, bookId); err != nil {
		return fmt.Errorf("failed to delete rating: %w", err)
	}
	
	return nil
}

func (r *PostgresBookRepository) GetRatings(bookId uuid.UUID, userId uuid.UUID) (*models.Rating, error) {
	var ratings *models.Rating
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

