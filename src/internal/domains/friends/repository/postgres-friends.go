package repository

import (

	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"

	"github.com/betterreads/internal/domains/friends/models"
    "github.com/google/uuid"
)


type PostgresFriendsRepository struct{
    db *sqlx.DB
}

func NewPostgresFriendsRepository (db *sqlx.DB) (FriendsRepository, error){
    schema := `
        CREATE TABLE IF NOT EXISTS friends (
            user_id UUID NOT NULL,
            friend_id UUID NOT NULL,
            friend_username VARCHAR(50) NOT NULL,
            user_username VARCHAR(50) NOT NULL,
            PRIMARY KEY (user_id, friend_id),
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (friend_id) REFERENCES users(id)
        );
    `

    if _, err := db.Exec(schema); err != nil {
        return nil, fmt.Errorf("failed to create friends table: %w", err)
    }

    schema = `
        CREATE TABLE IF NOT EXISTS friends_requests (
            user_id UUID NOT NULL,
            friend_id UUID NOT NULL,
            friend_username VARCHAR(50) NOT NULL,
            user_username VARCHAR(50) NOT NULL,
            PRIMARY KEY (user_id, friend_id),   
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (friend_id) REFERENCES users(id)
        );
    `
    if _, err := db.Exec(schema); err != nil {
        return nil, fmt.Errorf("failed to create friend_requests table: %w", err)
    }

    return &PostgresFriendsRepository{db: db}, nil
}

func (c PostgresFriendsRepository) GetFriends(userID uuid.UUID) ([]models.FriendOfUser, error){
    friends := []models.FriendOfUser{}
    query := `
        SELECT 
            CASE 
                WHEN user_id = $1 THEN friend_id 
                ELSE user_id 
            END as id,
            CASE 
                WHEN user_id = $1 THEN friend_username
                ELSE user_username 
            END as username
        FROM friends
        WHERE user_id = $1 OR friend_id = $1
    `
    err := c.db.Select(&friends, query, userID)
    if err != nil && err != sql.ErrNoRows {
        return nil, fmt.Errorf("failed to get friends: %w", err)
    }

    return friends, nil
}

func (c PostgresFriendsRepository)   AddFriend(userID uuid.UUID, friendID uuid.UUID) error {
    query:= `SELECT username FROM users WHERE id = $1`
    var friendUsername string
    err := c.db.Get(&friendUsername, query, friendID)
    if err != nil {
        return fmt.Errorf("failed to add friend: %w", err)
    }

    var userUsername string
    query = `SELECT username FROM users WHERE id = $1`
    err = c.db.Get(&userUsername, query, userID)
    if err != nil {
        return fmt.Errorf("failed to add friend: %w", err)
    }


    query = `INSERT INTO friends_requests (user_id, friend_id, friend_username, user_username) VALUES ($1, $2, $3, $4)`
    _, err = c.db.Exec(query, userID, friendID, friendUsername, userUsername)
    if err != nil {
        return fmt.Errorf("failed to add friend: %w", err)
    }
    return nil
}

func (c PostgresFriendsRepository)  AcceptFriendRequest(userID uuid.UUID, friendID uuid.UUID) error {
    type Usernames struct{
        FriendUsername string `db:"friend_username"`
        UserUsername string `db:"user_username"`
    }
    query := `SELECT friend_username, user_username FROM friends_requests WHERE user_id = $1 AND friend_id = $2`
    var usernames Usernames
    err := c.db.Get(&usernames, query, userID, friendID)
    if err != nil {
        return fmt.Errorf("failed to accept friend request: %w", err)
    }

    query = `DELETE FROM friends_requests WHERE user_id = $1 AND friend_id = $2`

    _, err = c.db.Exec(query, userID, friendID)
    if err != nil {
        return fmt.Errorf("failed to accept friend request: %w", err)
    }
    query = `INSERT INTO friends (user_id, friend_id, friend_username, user_username) VALUES ($1, $2, $3, $4)`
    _, err = c.db.Exec(query, userID, friendID, usernames.FriendUsername, usernames.UserUsername)
    if err != nil {
        return fmt.Errorf("failed to accept friend request: %w", err)
    }
    return nil
}

func (c PostgresFriendsRepository)  RejectFriendRequest(userID uuid.UUID, friendID uuid.UUID) error {
    query := `DELETE FROM friends_requests WHERE user_id = $1 AND friend_id = $2`
    _, err := c.db.Exec(query, userID, friendID)
    if err != nil && err != sql.ErrNoRows {
        return fmt.Errorf("failed to reject friend request: %w", err)
    }
    return nil
}


func (c PostgresFriendsRepository) GetFriendRequestsSent(senderId uuid.UUID) ([]models.FriendOfUser, error){
    query := `SELECT friend_id as id, friend_username as username FROM friends_requests WHERE user_id= $1`
    friends := []models.FriendOfUser{}
    err := c.db.Select(&friends, query, senderId)
    if err != nil && err != sql.ErrNoRows {
        return nil, fmt.Errorf("failed to get friend requests sent: %w", err)
    }
    return friends , nil
}

func (c PostgresFriendsRepository) GetFriendRequestsReceived(receiverId uuid.UUID) ([]models.FriendOfUser, error){
    query := `SELECT user_id as id, user_username as username FROM friends_requests WHERE friend_id = $1`
    friends := []models.FriendOfUser{}
    err := c.db.Select(&friends, query, receiverId)
    if err != nil && err != sql.ErrNoRows {
        return nil, fmt.Errorf("failed to get friend requests received: %w", err)
    }
    return friends , nil
}

func (c PostgresFriendsRepository)   CheckIfFriendRequestExists(userID uuid.UUID, friendID uuid.UUID) bool{
    query := `SELECT EXISTS (SELECT 1 FROM friends_requests WHERE user_id = $1 AND friend_id = $2)`
    var exists1 bool
    err := c.db.Get(&exists1, query, userID, friendID)
    if err != nil {
        return false
    }
    
    query = `SELECT EXISTS (SELECT 1 FROM friends_requests WHERE user_id = $1 AND friend_id = $2)`
    var exists2 bool
    err = c.db.Get(&exists2, query, friendID, userID)
    if err != nil {
        return false
    }
    return exists1 || exists2
}

func (c PostgresFriendsRepository) CheckIfFriendShipExists(userID uuid.UUID, friendID uuid.UUID) bool{
    query := `SELECT EXISTS (SELECT 1 FROM friends WHERE user_id = $1 AND friend_id = $2)`
    var exists1 bool
    err := c.db.Get(&exists1, query, userID, friendID)
    if err != nil {
        return false
    }
    
    query = `SELECT EXISTS (SELECT 1 FROM friends WHERE user_id = $1 AND friend_id = $2)`
    var exists2 bool
    err = c.db.Get(&exists2, query, friendID, userID)
    if err != nil {
        return false
    }
    return exists1 || exists2
}
