package repository

import (
	"database/sql"

	"github.com/betterreads/internal/domains/feed/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresFeedRepository struct {
	db *sqlx.DB
}

func NewPostgresFeedRepository(db *sqlx.DB) FeedRepository {
	return &PostgresFeedRepository{db: db}
}

func (pfr *PostgresFeedRepository) GetFeed(userId uuid.UUID) ([]models.Post, error) {
	posts := make([]models.Post, 0)

	mega_query := `
    select us.id as user_id, 
        us.username,  
        bk.id as book_id, 
        us.username as author_name, 
        bk.title, 
        bk.description, 
        bk.publication_date, 
        null AS rating
    from users us
    join friends fr on us.id = fr.user_a_id or us.id = fr.user_b_id 
    join books bk  on us.id = bk.author 
    where us.is_author = true 
        and ( fr.user_a_id = $1 or fr.user_b_id  = $1) 
        and us.id != $1
    union 
    select 
        us.id as user_id,
        us.username,
        bk.id as book_id,
        (select us.username from users us where us.id = bk.author) as author_name,
        bk.title,
        bk.description,
        r.publication_date,
        r.rating 
    from users us
    join friends fr on us.id = fr.user_a_id or us.id = fr.user_b_id 
    join reviews r on r.user_id = us.id
    join books bk on r.book_id = bk.id 
    where ( fr.user_a_id = $1 or fr.user_b_id  = $1) 
        and us.id != $1
    ORDER BY publication_date DESC;
    `
	err := pfr.db.Select(&posts, mega_query, userId)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	return posts, nil
}
