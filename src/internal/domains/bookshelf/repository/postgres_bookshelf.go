package repository

import (
	"database/sql"

	"fmt"

	"github.com/betterreads/internal/domains/bookshelf/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresBookShelfRepository struct {
	c *sqlx.DB
}

func NewPostgresBookShelfRepository(c *sqlx.DB) (BookshelfDatabase, error){
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
    bookShelfs, err := p.getBookShelfData(userId , status)
    if err != nil {
        return nil, err
    }

    reviews, err := p.getBookshelfReviews(userId, status)
    if err != nil {
        return nil, err
    }

    ratingStats, err := p.getBookshelfRatingStats(userId, status)
    if err != nil {
        return nil , err
    }


    resMap := map[uuid.UUID]*models.BookInShelfResponse{}
    for _, bookShelf:= range bookShelfs{
        bookShelf := &models.BookInShelfResponse{
            Date: bookShelf.Date,
            Status: bookShelf.Status,
            Title: bookShelf.Title,
            AuthorId: bookShelf.AuthorId,
            AuthorName: bookShelf.AuthorName,
            BookId: bookShelf.BookId,
        }
        resMap[bookShelf.BookId] = bookShelf
    }

    for _, reviewMap := range reviews{
        bookShelf := resMap[reviewMap.BookId]
        bookShelf.UserReview = reviewMap.Review
        bookShelf.UserRating = reviewMap.Rating
    }

    for _, ratingMap := range ratingStats{
        bookShelf := resMap[ratingMap.BookId]
        bookShelf.AvgRating = ratingMap.AvgRating
        bookShelf.TotalRatings = ratingMap.TotalRatings
    }


    res := make([]*models.BookInShelfResponse, 0, len(resMap))
    for _, bookShelf := range resMap {
        res = append(res, bookShelf)
    }

    return res , nil
}


func (p *PostgresBookShelfRepository) getBookShelfData(userId uuid.UUID, shelfType *models.BookShelfType) ([]models.BookShelfBasicData, error){
    bookShelfData:= []models.BookShelfBasicData{}

    query := `SELECT bk.title, bk.author as author_id, a.username as author_name, bk.id as book_id, bs.status, bs.date
                     FROM bookshelf bs
                     JOIN books bk ON bs.book_id=bk.id
                     JOIN users a ON bk.author=a.id
                     WHERE bs.user_id=$1 AND ($2::VARCHAR IS NULL OR bs.status=$2);`

    err := p.c.Select(&bookShelfData, query, userId, shelfType)
    if err != nil && err != sql.ErrNoRows {
        return nil, fmt.Errorf("failed to get bookshelf: %w", err)
    } 

    return bookShelfData, nil

}

func (p *PostgresBookShelfRepository) getBookshelfReviews(userId uuid.UUID, shelfType *models.BookShelfType) ([]models.BookShelfUserReview, error){
    reviewsMaps := []models.BookShelfUserReview{}
    query := `SELECT r.book_id, r.review, r.rating
             FROM reviews r
             WHERE r.user_id = $1 AND r.book_id IN (
                SELECT bs.book_id 
                FROM bookshelf bs 
                WHERE bs.user_id=$1 AND ($2::VARCHAR IS NULL OR bs.status=$2)
             );`
    err := p.c.Select(&reviewsMaps, query, userId, shelfType)
    if err != nil && err != sql.ErrNoRows {
        return nil, fmt.Errorf("failed to get reviews: %w", err)
    } 
    return reviewsMaps, nil
}

func (p *PostgresBookShelfRepository) getBookshelfRatingStats(userId uuid.UUID, shelfType *models.BookShelfType) ([]models.BookShelfRatingStats, error){
    RatingStatsMaps := []models.BookShelfRatingStats{}
    query := `SELECT r.book_id, COALESCE(AVG(r.rating),0) as avg_ratings, COUNT(r.rating) as total_ratings
                FROM reviews r
                WHERE r.user_id = $1 AND r.book_id IN (
                    SELECT book_id 
                    FROM bookshelf bs
                    WHERE user_id=$1 AND ($2::VARCHAR IS NULL OR bs.status=$2))
                GROUP BY r.book_id;`


    err := p.c.Select(&RatingStatsMaps, query, userId, shelfType)
        if err != nil && err != sql.ErrNoRows {
            return nil, fmt.Errorf("failed to get rating stats: %w", err)
        }
    return RatingStatsMaps, nil
}

func (p *PostgresBookShelfRepository) AddBookToShelf(userId uuid.UUID, req *models.BookShelfRequest) error {
    query  := `INSERT INTO bookshelf (user_id, book_id, status, date)
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

