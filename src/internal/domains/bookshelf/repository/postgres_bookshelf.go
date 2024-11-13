package repository

import (
	"database/sql"

	"fmt"

	"github.com/betterreads/internal/domains/bookshelf/models"
	booksRepo "github.com/betterreads/internal/domains/books/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresBookShelfRepository struct {
	c *sqlx.DB
}

func NewPostgresBookShelfRepository(c *sqlx.DB) (BookshelfDatabase, error) {
	enableUUIDExtension := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	if _, err := c.Exec(enableUUIDExtension); err != nil {
		return nil, fmt.Errorf("failed to enable uuid extension: %w", err)
	}

	schemaBookshelf := `
        CREATE TABLE IF NOT EXISTS bookshelf (
            user_id UUID NOT NULL, 
            book_id UUID NOT NULL, 
            status VARCHAR(50) NOT NULL,
            date DATE NOT NULL,
            PRIMARY KEY (user_id, book_iD),
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (book_id) REFERENCES books(id)
        );
    `

	if _, err := c.Exec(schemaBookshelf); err != nil {
		return nil, fmt.Errorf("failed to create bookshelf table: %w", err)
	}

	return &PostgresBookShelfRepository{c: c}, nil
}

func (p *PostgresBookShelfRepository) GetBookShelf(userId uuid.UUID, shelfType models.BookShelfType) ([]*models.BookInShelfResponse, error) {
	var status *models.BookShelfType

	if shelfType == models.BookShelfAll {
		status = nil
	} else {
		status = &shelfType
	}
    
    res:= []*models.BookInShelfResponse{}
    query := `
    WITH ratings AS (
        SELECT
            r.book_id,
            COALESCE(AVG(r.rating),0) as avg_ratings,
            COUNT(*) as total_ratings
        FROM reviews r
        group by r.book_id
    ),
    user_ratings AS (
        SELECT 
            r.book_id,
            r.review,
            rating
        FROM reviews r
        WHERE r.user_id= $1
    )
    SELECT 
        bk.title,
        bk.author as author_id,
        (SELECT username FROM users WHERE id=bk.author) as author_name,
        bk.id as book_id,
        bs.status,
        bs.date,
        COALESCE(r.avg_ratings, 0) as avg_ratings,
        COALESCE(r.total_ratings,0) as total_ratings,
        COALESCE(ur.review, '') as user_review,
        COALESCE(ur.rating, 0) as user_rating
    FROM bookshelf bs
    JOIN books bk ON bs.book_id=bk.id
    LEFT JOIN ratings r ON r.book_id=bk.id
    LEFT JOIN user_ratings ur ON ur.book_id=bk.id
    WHERE bs.user_id=$1 AND ($2::VARCHAR IS NULL OR bs.status=$2)
    ORDER BY avg_ratings DESC;
    `
    if err := p.c.Select(&res, query, userId, status); err != nil {
        if err != sql.ErrNoRows {
            return nil, fmt.Errorf("failed to get bookshelf: %w", err)
        }
    }
    
    return res, nil
}


func (p *PostgresBookShelfRepository) AddBookToShelf(userId uuid.UUID, req *models.BookShelfRequest) error {
	query := `INSERT INTO bookshelf (user_id, book_id, status, date)
                      VALUES ($1, $2, $3, now())`

	_, err := p.c.Exec(query, userId, req.BookId, req.Status)
	if err != nil {
		return fmt.Errorf("failed to add book to shelf: %w", err)
	}

	return nil
}

func (p *PostgresBookShelfRepository) EditBookInShelf(userId uuid.UUID, req *models.BookShelfRequest) error {
	query := ` UPDATE bookshelf
                      SET status=$1, date=now()
                      WHERE user_id=$2 AND book_id=$3;`
	_, err := p.c.Exec(query, req.Status, userId, req.BookId)
	if err != nil {
		return fmt.Errorf("failed to edit book in shelf: %w", err)
	}

	return nil
}

func (p *PostgresBookShelfRepository) CheckIfBookIsInUserShelf(userId uuid.UUID, bookId uuid.UUID) bool {
	exists := false
	query := `SELECT EXISTS(SELECT 1 FROM bookshelf WHERE user_id=$1 AND book_id=$2);`
	err := p.c.Get(&exists, query, userId, bookId)
	if err != nil {
		return false
	}

	return exists
}

func (p *PostgresBookShelfRepository) DeleteBookFromShelf(userId uuid.UUID, bookId uuid.UUID) error {
	query := `DELETE FROM bookshelf WHERE user_id=$1 AND book_id=$2;`
	_, err := p.c.Exec(query, userId, bookId)
	if err != nil {
		return fmt.Errorf("failed to delete book from shelf: %w", err)
	}

	return nil
}

func (p *PostgresBookShelfRepository) SearchBookShelf(userId uuid.UUID, shelfType models.BookShelfType, genre string, sort string, isDirAsc bool) ([]*models.BookInShelfResponse, error) {
	var status *models.BookShelfType
	if shelfType == models.BookShelfAll {
		status = nil
	} else {
		status = &shelfType
	}

    books := []*models.BookInShelfResponse{}

    query := `
    WITH ratings AS (
        SELECT
            r.book_id,
            COALESCE(AVG(r.rating),0) as avg_ratings,
            COUNT(*) as total_ratings
        FROM reviews r
        group by r.book_id
    ),
    user_ratings AS (
        SELECT 
            r.book_id,
            r.review,
            rating
        FROM reviews r
        WHERE r.user_id= $1
    )
    SELECT 
        bk.title,
        bk.author as author_id,
        (SELECT username FROM users WHERE id=bk.author) as author_name,
        bk.id as book_id,
        bs.status,
        bs.date,
        COALESCE(r.avg_ratings, 0) as avg_ratings,
        COALESCE(r.total_ratings,0) as total_ratings,
        COALESCE(ur.review, '') as user_review,
        COALESCE(ur.rating, 0) as user_rating
    FROM bookshelf bs
    JOIN books bk ON bs.book_id=bk.id
    LEFT JOIN ratings r ON r.book_id=bk.id
    LEFT JOIN user_ratings ur ON ur.book_id=bk.id
    `
    var genreId int
    var err error
    if genre != "" {
        genreId, err = booksRepo.GetGenreById(genre)
        if err != nil {
            return nil, fmt.Errorf("failed to get genre id: %w", err)
        }
        query += "JOIN genres_books bg ON bg.book_id=bk.id"
    }

    query += " WHERE bs.user_id=$1 AND ($2::VARCHAR IS NULL OR bs.status=$2)"

    if genre != "" {
        query += " AND bg.genre_id=$3"
    }

	if sort != "" {
		var direciton string
		if isDirAsc{
			direciton = "ASC"
		} else {
			direciton = "DESC"
		}
		query += " ORDER BY " + sort + " " + direciton
	}

    if genre != "" {
        if  err := p.c.Select(&books, query, userId, status, genreId); err != nil {
            return nil, fmt.Errorf("failed to search books in shelf: %w", err)
        }
    } else {
        if  err := p.c.Select(&books, query, userId, status); err != nil {
            return nil, fmt.Errorf("failed to search books in shelf: %w", err)
        }
    }

    return books, nil
}
