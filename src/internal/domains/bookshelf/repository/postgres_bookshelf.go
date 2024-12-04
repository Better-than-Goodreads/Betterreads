package repository

import (
	"database/sql"
	"strconv"
	"strings"

	"fmt"

	booksRepo "github.com/betterreads/internal/domains/books/repository"
	"github.com/betterreads/internal/domains/bookshelf/models"
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

const query_beggining = `
   WITH ratings AS (
        SELECT
            r.book_id,
            COALESCE(AVG(r.rating),0) as avg_ratings,
            COUNT(*) as total_ratings
        FROM reviews r
        GROUP BY r.book_id
    ),
    user_ratings AS (
        SELECT 
            r.book_id,
            r.review,
            r.rating
        FROM reviews r
        WHERE r.user_id = $1
    )
    SELECT 
        bk.title,
        bk.author as author_id,
		u.username as author_name,
        bk.description, 
        bk.publication_date,
        bs.date,
        bk.language,
        array_agg(bg.genre_id) as genres,
        bk.amount_of_pages,
        COALESCE(r.total_ratings, 0) as total_ratings,
        COALESCE(r.avg_ratings, 0) as avg_ratings,
        bs.status,
        COALESCE(ur.review, '') as user_review,
        COALESCE(ur.rating, 0) as user_rating,
        bk.id as id
    FROM bookshelf bs
    JOIN books bk ON bs.book_id = bk.id
	JOIN users u ON bk.author = u.id
    LEFT JOIN ratings r ON r.book_id = bk.id
    LEFT JOIN user_ratings ur ON ur.book_id = bk.id
    LEFT JOIN genres_books bg ON bg.book_id = bk.id 
	WHERE
`

const query_normal_cond = `
	($2::VARCHAR IS NULL OR bs.status=$2) AND bs.user_id=$1
`

const query_group_by = `
    GROUP BY
	bk.id,                 
    bk.title,
    bk.author,
    u.username,
    bk.description,
    bk.publication_date,
    bs.date,
    bk.language,
    bk.amount_of_pages,
    bs.status,
	total_ratings,
	avg_ratings,
	ur.review,
	ur.rating
`

func (p *PostgresBookShelfRepository) GetBookShelf(userId uuid.UUID, shelfType models.BookShelfType) ([]*models.BookInShelfResponse, error) {
	var status *models.BookShelfType

	if shelfType == models.BookShelfAll {
		status = nil
	} else {
		status = &shelfType
	}

	books := []*models.BookInShelfResponse{}
	query := query_beggining + query_normal_cond + query_group_by + `
	ORDER BY bs.date DESC;
    `

	if err := p.c.Select(&books, query, userId, status); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to get bookshelf: %w", err)
		}
	}

	for _, book := range books {
		book.GenresArray = parseGenres(book.Genres)
		book.Genres = ""
	}

	return books, nil

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

	query := query_beggining

	var args []interface{}
	args = append(args, userId, status)
	var genreId int
	var err error
	if genre != "" {
		genreId, err = booksRepo.GetGenreById(genre)
		if err != nil {
			return nil, fmt.Errorf("failed to get genre id: %w", err)
		}
		query += "exists(SELECT 1 FROM genres_books WHERE book_id=bk.id AND genre_id=$3) AND "
		args = append(args, genreId)
	}

	query += query_normal_cond + query_group_by

	if sort != "" {
		var direciton string
		if isDirAsc {
			direciton = "ASC"
		} else {
			direciton = "DESC"
		}
		query += " ORDER BY " + sort + " " + direciton
	}

	if err := p.c.Select(&books, query, args...); err != nil {
		return nil, fmt.Errorf("failed to search books in shelf: %w", err)
	}

	for _, book := range books {
		book.GenresArray = parseGenres(book.Genres)
		book.Genres = ""
	}

	return books, nil
}

func parseGenres(genres string) *[]string {
	// Remove the curly braces
	genres = strings.Trim(genres, "{}")
	// Split the string by commas
	genresArr := strings.Split(genres, ",")
	res := make([]string, 0, len(genresArr))

	for _, genre := range genresArr {
		genreId, _ := strconv.Atoi(genre)
		genre_str := booksRepo.GetGenre(genreId)
		res = append(res, genre_str)
	}
	return &res
}
