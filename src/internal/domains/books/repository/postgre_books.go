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
	return -1, ErrGenreNotFound
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
			author UUID  NOT NULL,
			description VARCHAR(255) NOT NULL,
			amount_of_pages INTEGER NOT NULL,
			publication_date VARCHAR(255) NOT NULL,
            language VARCHAR(255) NOT NULL,
			FOREIGN KEY (author) REFERENCES users(id)
			);
			`
	schemaGendersBooks := `
        CREATE TABLE IF NOT EXISTS genres_books (
            book_id UUID NOT NULL,
            genre_id INT,
            PRIMARY KEY (book_id, genre_id),
            FOREIGN KEY (book_id) REFERENCES books(id)
        );
    `
    
    schemaReviews := `
        CREATE TABLE IF NOT EXISTS reviews (
            user_id UUID,
            book_id UUID,
            rating INTEGER,
            review VARCHAR (255),
            PRIMARY KEY (user_id, book_id),
            FOREIGN KEY (book_id) REFERENCES books(id),
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
    `
		
	schemaPictures := `
		CREATE TABLE IF NOT EXISTS pictures (
			book_id UUID,
			picture BYTEA,
			FOREIGN KEY (book_id) REFERENCES books(id),
			PRIMARY KEY (book_id)
		);
		`
	
	if _, err := c.Exec(schemaBooks); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	if _, err := c.Exec(schemaReviews); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	if _, err := c.Exec(schemaGendersBooks); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	if _, err := c.Exec(schemaPictures); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PostgresBookRepository{c}, nil
}

func (r *PostgresBookRepository) SaveBook(book *models.NewBookRequest, author uuid.UUID) (*models.Book, error) {
	bookRecord := &models.BookDb{}
	query := `INSERT INTO books (title, author, description,  amount_of_pages,
                    publication_date, language)
                    VALUES ($1, $2, $3, $4, $5, $6)
                    RETURNING id, title, author, description, amount_of_pages, publication_date, language;`

	args := []interface{}{book.Title, author, book.Description, book.AmountOfPages, book.PublicationDate, book.Language}

	if err := r.c.Get(bookRecord, query, args...); err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
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

	query = `INSERT INTO pictures (book_id, picture)
			VAlUES ($1, $2);`

	args = []interface{}{bookRecord.Id, book.Picture}

	if _, err := r.c.Exec(query, args...); err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
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

func (r *PostgresBookRepository) getRatingsForBook(book_id uuid.UUID) (*models.Ratings, error) {
	ratings := &models.Ratings{}
	query := `SELECT COUNT(*) as total_ratings, COALESCE(AVG(rating), 0) as avg_ratings FROM reviews WHERE book_id = $1;`
	if err := r.c.Get(ratings, query, book_id); err != nil {
		if err == sql.ErrNoRows {
			return ratings, nil
		}
		return nil, fmt.Errorf("failed to get ratings: %w", err)
	}
	return ratings, nil
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
    book, err := r.getBookInfo(bookdb)
    if err != nil {
        return nil, fmt.Errorf("failed to get book: %w", err)
    }
	return book, nil
}

func (r *PostgresBookRepository) GetBooksByName(name string) ([]*models.Book, error) {
    books := &[]*models.BookDb{}
    query := `SELECT * FROM books WHERE LOWER(title) LIKE LOWER('%'||$1||'%');`
    if err := r.c.Select(books, query, name); err != nil {
        return nil, fmt.Errorf("failed to get book: %w", err)
    }
    res := []*models.Book{}
    for _, book := range *books {
        Bookres, err := r.getBookInfo(book)
        if err != nil {
            return nil, fmt.Errorf("failed to get book: %w", err)
        }
        res = append(res, Bookres)
    }
    return res, nil
}

func (r *PostgresBookRepository) GetBookPictureById(id uuid.UUID) ([]byte, error) {
	var picture []byte
	query := `SELECT picture FROM pictures WHERE book_id = $1;`
	if err := r.c.Get(&picture, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBookNotFound
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	return picture, nil
}

func (r *PostgresBookRepository) GetBooks() ([]*models.Book, error) {
	var books []*models.BookDb
	query := `SELECT * FROM books;`
	if err := r.c.Select(&books, query); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no books found")
		}
		return nil, fmt.Errorf("failed to get books: %w", err)
	}
	res := []*models.Book{}
    for _, book := range books {
        Bookres, err := r.getBookInfo(book)
        if err != nil {
            return nil, fmt.Errorf("failed to get book: %w", err)
        }
        res = append(res, Bookres)
    }
	return res, nil
}

// func (r *PostgresBookRepository) RateBook(bookId uuid.UUID, userId uuid.UUID, rating int) error {
// 	exists, err := checkIfRatingExists(r.c, bookId, userId)
// 	if err != nil {
// 		return fmt.Errorf("failed to check rating: %w", err)
// 	}
//
// 	if !exists {
// 		query := `INSERT INTO ratings (user_id, book_id, rating)
// 				  VALUES ($1, $2, $3);`
// 		args := []interface{}{userId, bookId, rating}
//
// 		_, err := r.c.Exec(query, args...)
// 		if err != nil {
// 			return fmt.Errorf("failed to rate book: %w", err)
// 		}
//
// 		return nil
// 	} else {
// 		query := `UPDATE ratings SET rating = $1 WHERE user_id = $2 AND book_id = $3;`
// 		args := []interface{}{rating, userId, bookId}
// 		_, err := r.c.Exec(query, args...)
// 		if err != nil {
// 			return fmt.Errorf("failed to rate book: %w", err)
// 		}
// 		return nil
// 	}
// }

func checkIfReviewExists(c *sqlx.DB, bookId uuid.UUID, userId uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM reviews WHERE user_id = $1 AND book_id = $2);`
	var exists bool
	if err := c.Get(&exists, query, userId, bookId); err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

// func (r *PostgresBookRepository) DeleteRating(bookId uuid.UUID, userId uuid.UUID) error {
// 	exists, err := checkIfRatingExists(r.c, bookId, userId)
// 	if err != nil {
// 		return fmt.Errorf("failed to check rating: %w", err)
// 	}
//
// 	if exists {
// 		query := `DELETE FROM ratings WHERE user_id = $1 AND book_id = $2;`
//
// 		if _, err := r.c.Exec(query, userId, bookId); err != nil {
// 			return fmt.Errorf("failed to delete rating: %w", err)
// 		}
//
// 		return nil
// 	} else {
// 		return ErrRatingNotFound
// 	}
// }

func (r *PostgresBookRepository) GetRatingUser(bookId uuid.UUID, userId uuid.UUID) (*models.Rating, error) {
	var ratings models.Rating
	query := `SELECT * FROM reviews WHERE book_id = $1 AND user_id = $2;`
	args := []interface{}{bookId, userId}

	if err := r.c.Get(&ratings, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRatingNotFound
		}
		return nil, fmt.Errorf("failed to get ratings: %w", err)
	}
	return &ratings, nil
}

func (r *PostgresBookRepository) GetAuthorName(authorId uuid.UUID) (string, error) {
	var authorName string
	query := `SELECT username FROM users WHERE id = $1;`
	if err := r.c.Get(&authorName, query, authorId); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("author with id %s not found", authorId)
		}
		return "", fmt.Errorf("failed to get author: %w", err)
	}
	return authorName, nil
}


func (r *PostgresBookRepository) AddReview(bookId uuid.UUID, userId uuid.UUID, review string, rating int) error {
    query := `SELECT EXISTS(SELECT 1 FROM reviews WHERE user_id = $1 and book_id = $2);`
    args := []interface{}{userId, bookId}
    var exists bool
    if err := r.c.Get(&exists, query, args...); err != nil {
        return fmt.Errorf("failed to add review: %w", err)
    }

    query = `INSERT INTO reviews (user_id, book_id, review, rating)
    VALUES ($1, $2, $3, $4);`
    args = []interface{}{userId, bookId, review, rating}

    if _, err := r.c.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to add review: %w", err)
    }
    return nil
}

func (r *PostgresBookRepository) DeleteReview(bookId uuid.UUID, userId uuid.UUID) error {
    return nil
}


func (r *PostgresBookRepository) getBookInfo(book *models.BookDb) (*models.Book, error) {
    genres, err := r.getGenresForBook(book.Id)
    if err != nil {
        return nil, fmt.Errorf("failed to get books: %w", err)
    }
    ratings, err := r.getRatingsForBook(book.Id)
    if err != nil {
        return nil, fmt.Errorf("failed to get book: %w", err)
    }
    return utils.MapBookDbToBook(book, genres, ratings), nil
}
