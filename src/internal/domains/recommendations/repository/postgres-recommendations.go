package repository

import (
	"database/sql"
	"fmt"
	"sort"

	bm "github.com/betterreads/internal/domains/books/models"
	"github.com/betterreads/internal/domains/books/repository"
	bsm "github.com/betterreads/internal/domains/bookshelf/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresRecommendationsRepository struct {
	br repository.BooksDatabase
	c  *sqlx.DB
}

func NewPostgresRecommendationsRepository(c *sqlx.DB, br repository.BooksDatabase) RecommendationsDatabase {
	return &PostgresRecommendationsRepository{c: c, br: br}
}

func (r *PostgresRecommendationsRepository) GetRecommendations(userId uuid.UUID) (map[string][]*bm.Book, error) {
	preferedGenres, err := r.getPreferedGenres(userId)
	if err != nil {
		return nil, err
	}

	booksByGenre := make(map[string][]*bm.Book)
	for _, genre := range preferedGenres {
		fmt.Printf("Getting books for genre: %v \n", genre)
		books, err := r.GetPreferedBooks(genre, 5, userId)
		if err != nil {
			return nil, fmt.Errorf("failed to get books by genre: %w", err)
		}
		booksByGenre[genre] = books
	}
	return booksByGenre, nil
}

func (r *PostgresRecommendationsRepository) GetMoreRecommendations(userId uuid.UUID, genre string) ([]*bm.Book, error) {
	books, err := r.GetPreferedBooks(genre, 20, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get books by genre: %w", err)
	}
	return books, nil
}

func (r *PostgresRecommendationsRepository) GetPreferedBooks(genre string, limit int, userId uuid.UUID) ([]*bm.Book, error) {
	genre_id, err := getGenreId(genre)
	if err != nil {
		return nil, fmt.Errorf("failed to get genre id: %w", err)
	}

	// Gets books by genre that user has not read
	query := `SELECT bk.title, bk.author, bk.description, bk.amount_of_pages, bk.
              publication_date, bk.language, bk.id 
              FROM books bk
              JOIN genres_books gb ON bk.id = gb.book_id
              WHERE gb.genre_id = $1 AND bk.id NOT IN (SELECT book_id FROM bookshelf WHERE user_id = $2)
              LIMIT $3;`

	books := []*bm.BookDb{}
	err = r.c.Select(&books, query, genre_id, userId, limit)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get preferedBooks: %w", err)
	}

	res := []*bm.Book{}
	for _, book := range books {
		bookRes, err := r.br.GetBookInfo(book)
		if err != nil {
			return nil, fmt.Errorf("failed to get book info: %w", err)
		}
		res = append(res, bookRes)
	}
    
    sort.Slice(res, func(i, j int) bool {
        return res[i].AverageRating > res[j].AverageRating
    })

	return res, nil
}

func (r *PostgresRecommendationsRepository) getPreferedGenres(userId uuid.UUID) ([]string, error) {
	// First gets the books that user has read from bookshelf
	query := `SELECT bk.title, bk.author, bk.description, bk.amount_of_pages, bk.
              publication_date, bk.language, bk.id 
              FROM bookshelf bs
              JOIN books bk ON bs.book_id = bk.id
              WHERE bs.user_id = $1 AND bs.status='read';`
	userBooks := []bm.BookDb{}
	err := r.c.Select(&userBooks, query, userId)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get user books: %w", err)
	}

	genresMap := make(map[string]int)
	// Then gets the genres of the books
	for _, book := range userBooks {
		genres, err := r.br.GetGenresForBook(book.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get genres for book: %w", err)
		}
		for _, genre := range genres {
			genresMap[genre]++
		}
	}

	// Sorts the genres by the amount of times they appear
	topGenres := r.sortGenres(genresMap)

	return topGenres, nil
}

func (r *PostgresRecommendationsRepository) sortGenres(genresMap map[string]int) []string {
	type GenreCount struct {
		Genre string
		Count int
	}

	var genres []GenreCount
	for genre, count := range genresMap {
		genres = append(genres, GenreCount{Genre: genre, Count: count})
	}

	sort.Slice(genres, func(i, j int) bool {
		return genres[i].Count > genres[j].Count
	})

	topGenres := []string{}
	for i := 0; i < len(genres) && i < 3; i++ {
		topGenres = append(topGenres, genres[i].Genre)
	}

	return topGenres
}

func getGenreId(genre string) (int, error) {
	for key, value := range repository.GenresDict {
		if value == genre {
			return key, nil
		}
	}
	return -1, repository.ErrGenreNotFound
}

func (r *PostgresRecommendationsRepository) CheckIfUserHasValidShelf(userId uuid.UUID) bool {
	bookshelf := []bsm.BookShelfBasicData{}
	query := `SELECT * FROM bookshelf WHERE user_id=$1 AND status='read';`
	err := r.c.Select(&bookshelf, query, userId)
	if err != nil {
		return false
	}
	return len(bookshelf) >= 5
}
