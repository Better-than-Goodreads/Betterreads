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
            rating INTEGER NOT NULL,
            review VARCHAR (255) NOT NULL,
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
			return nil, ErrGenreNotFound
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

    authorName, err := r.getAuthorName(author)
    if err != nil {
        return nil, fmt.Errorf("failed to create book: %w", err)
    }

    res := utils.MapBookDbToBook(bookRecord, book.Genres, &models.Ratings{}, authorName)

	return res, nil
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
			return nil, ErrBookNotFound
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
        if err != sql.ErrNoRows {
            return []*models.Book{}, nil
        }
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
			return nil, nil
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
            return []*models.Book{}, nil
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
    fmt.Printf("Books: %v\n", res)
	return res, nil
}

func (r *PostgresBookRepository) GetBooksOfAuthor(authorId uuid.UUID) ([]*models.Book, error) {
	var books []*models.BookDb
    query := `SELECT * FROM books WHERE author = $1;`
	if err := r.c.Select(&books, query, authorId); err != nil {
		if err == sql.ErrNoRows {
			return []*models.Book{}, nil
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

func (r *PostgresBookRepository) RateBook(bookId uuid.UUID, userId uuid.UUID, rating int) (*models.Rating, error) {
	var ratingRecord models.Rating
	query := `INSERT INTO reviews (user_id, book_id, rating, review)
			VALUES ($1, $2, $3, $4)
			RETURNING user_id, book_id, rating;`
	args := []interface{}{userId, bookId, rating,""}

	if err := r.c.Get(&ratingRecord, query, args...); err != nil {
		return nil, fmt.Errorf("failed to rate book: %w", err)
	}
	return &ratingRecord, nil
}

func (r *PostgresBookRepository) CheckIfRatingExists(bookId uuid.UUID, userId uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM reviews WHERE user_id = $1 AND book_id = $2);`
	var exists bool
	if err := r.c.Get(&exists, query, userId, bookId); err != nil {
		return false, err
	} else {
		return exists, nil
	}
}

func (r *PostgresBookRepository) UpdateRating(bookId uuid.UUID, userId uuid.UUID, rating int) (error){
    query := `UPDATE reviews SET rating = $1 WHERE user_id = $2 AND book_id = $3;`
    args := []interface{}{rating, userId, bookId}
    if _, err := r.c.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update rating: %w", err)
    }

    return nil
}

func (r *PostgresBookRepository) GetBookReviewOfUser(bookId uuid.UUID, userId uuid.UUID) (*models.Review, error) {
	var ratings models.ReviewDb
	query := `SELECT * FROM reviews WHERE book_id = $1 AND user_id = $2;`
	args := []interface{}{bookId, userId}

	if err := r.c.Get(&ratings, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrReviewNotFound
		}
		return nil, fmt.Errorf("failed to get ratings: %w", err)
	}
    ReviewRes:= &models.Review{
        Text: ratings.Review,
        Rating: ratings.Rating,
    }
	return ReviewRes, nil
}

func (r *PostgresBookRepository) GetAllReviewsOfUser(userId uuid.UUID) ([]*models.ReviewOfUser, error) {
	var username string
	query := `SELECT username FROM users WHERE id = $1;`
	if err := r.c.Get(&username, query, userId); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get author: %w", err)
	}

    res := []*models.ReviewOfUser{}  
    query = `
        SELECT b.title AS book_title, r.review,b.id as book_id, r.rating
        FROM reviews r
        INNER JOIN books b ON r.book_id = b.id
        WHERE r.user_id = $1;
    `
    if err := r.c.Select(&res, query, userId); err != nil {
        if err == sql.ErrNoRows {
            return []*models.ReviewOfUser{}, nil
        }
        return nil, fmt.Errorf("failed to get reviews: %w", err)
    }

    return res, nil
}

func (r *PostgresBookRepository) getAuthorName(authorId uuid.UUID) (string, error) {
	var authorName string
	query := `SELECT username FROM users WHERE id = $1;`
	if err := r.c.Get(&authorName, query, authorId); err != nil {
		if err == sql.ErrNoRows {
			return "", ErrAuthorNotFound
		}
		return "", fmt.Errorf("failed to get author: %w", err)
	}
	return authorName, nil
}

func (r *PostgresBookRepository) AddReview(bookId uuid.UUID, userId uuid.UUID, review string, rating int) error {
    args := []interface{}{userId, bookId}
    query := `INSERT INTO reviews (user_id, book_id, review, rating)
    VALUES ($1, $2, $3, $4);`
    args = []interface{}{userId, bookId, review, rating}

    if _, err := r.c.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to add review: %w", err)
    }
    return nil
}

func (r *PostgresBookRepository) EditReview(bookId uuid.UUID, userId uuid.UUID, rating int, review string) error{
    fmt.Println("edit review")
    query := `UPDATE reviews SET review = $1, rating = $2 WHERE book_id = $3 AND user_id = $4;`
    args := []interface{}{review, rating, bookId, userId}
    if _, err := r.c.Exec(query, args...); err != nil {
        return fmt.Errorf("failed to update review: %w", err)
    }
    return nil
}

func (r *PostgresBookRepository) DeleteReview(bookId uuid.UUID, userId uuid.UUID) error {
    return nil //lol
}

func (r *PostgresBookRepository) CheckifReviewExists(bookId uuid.UUID, userId uuid.UUID) (bool, error) {
    reviewCheck := &models.ReviewDb{}
    query := `SELECT * FROM reviews WHERE book_id = $1 AND user_id = $2;`
    args := []interface{}{bookId, userId}
    if err := r.c.Get(reviewCheck, query, args...); err != nil {
        if err == sql.ErrNoRows{
            return false , nil
        }
        return false, fmt.Errorf("failed to get review: %w", err)
    }

    if reviewCheck.Review == "" {
        return false, ErrReviewEmpty
    }

    return true, nil
}

func (r *PostgresBookRepository) GetBookReviews(bookID uuid.UUID) ([]*models.ReviewOfBook, error){
    res := []*models.ReviewOfBook{}
    query := `
        SELECT u.username, r.review, u.id AS user_id, r.rating
        FROM reviews r
        INNER JOIN users u ON r.user_id = u.id
        WHERE r.book_id = $1;
    `

    if err := r.c.Select(&res, query, bookID); err != nil {
        if err == sql.ErrNoRows {
            return []*models.ReviewOfBook{}, nil
        }
        return nil, fmt.Errorf("failed to get reviews: %w", err)
    }

    return res, nil
}


func (r *PostgresBookRepository) getBookInfo(book *models.BookDb) (*models.Book, error) {
    genres, err := r.getGenresForBook(book.Id)
    if err != nil {
        return nil, fmt.Errorf("failed to get books: %w", err)
    }
    ratings, err := r.getRatingsForBook(book.Id)
    if err != nil {
        if err != ErrRatingNotFound {
            return nil, fmt.Errorf("failed to get book: %w", err)
        }
    }
    author, err := r.getAuthorName(book.Author)
    if err != nil {
        return nil, fmt.Errorf("failed to get book: %w", err)
    }

    return utils.MapBookDbToBook(book, genres, ratings, author), nil
}

func (r *PostgresBookRepository) CheckIfBookExists(bookId uuid.UUID) bool{
    exists := false
    query := `SELECT EXISTS(SELECT 1 FROM books WHERE id = $1);`
    if err := r.c.Get(&exists, query, bookId); err != nil {
        return false
    }
    return exists
}

func (r *PostgresBookRepository) CheckIfAuthorExists(authorId uuid.UUID) bool {
    exists := false
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);`
    if err := r.c.Get(&exists, query, authorId); err != nil {
        return false
    }
    return exists
}
