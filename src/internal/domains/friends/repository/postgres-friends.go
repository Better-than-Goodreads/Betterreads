package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
    
    um "github.com/betterreads/internal/domains/users/models"
    "github.com/betterreads/internal/domains/users/utils"
	"github.com/google/uuid"
)

type PostgresFriendsRepository struct {
	db *sqlx.DB
}

func NewPostgresFriendsRepository(db *sqlx.DB) (FriendsRepository, error) {
	schema := `
        CREATE TABLE IF NOT EXISTS friends (
            user_a_id UUID NOT NULL,
            user_b_id UUID NOT NULL,
            PRIMARY KEY (user_a_id, user_b_id),
            FOREIGN KEY (user_a_id) REFERENCES users(id),
            FOREIGN KEY (user_b_id) REFERENCES users(id)
        );
    `

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create friends table: %w", err)
	}

	schema = `
        CREATE TABLE IF NOT EXISTS friends_requests (
            recipient_id UUID NOT NULL,
            sender_id UUID NOT NULL,
            PRIMARY KEY (recipient_id, sender_id),   
            FOREIGN KEY (recipient_id) REFERENCES users(id),
            FOREIGN KEY (sender_id) REFERENCES users(id)
        );
    `
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create friend_requests table: %w", err)
	}

	return &PostgresFriendsRepository{db: db}, nil
}

func (c PostgresFriendsRepository) GetFriends(userID uuid.UUID) ([]um.UserResponse, error) {
	friends := []um.UserRecord{}
	query := `
        SELECT 
            us.* 
        FROM friends fr
        JOIN users us ON us.id = 
            CASE 
                WHEN fr.user_a_id = $1 THEN fr.user_b_id
                ELSE fr.user_a_id 
            END
        WHERE fr.user_a_id = $1 OR fr.user_b_id = $1
    `

	err := c.db.Select(&friends, query, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get friends: %w", err)
	}
    
    res := []um.UserResponse{}
    for _, friend := range friends {
        res = append(res, *utils.MapUserRecordToUserResponse(&friend))
    }

	return res, nil
}

func (c PostgresFriendsRepository) AddFriend(senderId uuid.UUID, recipientId uuid.UUID) error {
    query := `INSERT INTO friends_requests (sender_id, recipient_id) VALUES ($1, $2)`
    _, err := c.db.Exec(query, senderId, recipientId)
	if err != nil {
		return fmt.Errorf("failed to add friend: %w", err)
	}
	return nil
}

func (c PostgresFriendsRepository) AcceptFriendRequest(recipientId uuid.UUID, senderId uuid.UUID) error {
    err := c.DeleteRequest(senderId, recipientId)
    if err != nil {
        return fmt.Errorf("failed to accept friend request: %w", err)
    }
    
    // checks if the other user has also sent a friend request from the other side
    exits := c.CheckIfFriendRequestExists(recipientId, senderId)
    if exits {
        err = c.DeleteRequest(recipientId, senderId)
        if err != nil {
            return fmt.Errorf("failed to accept friend request: %w", err)
        }
    }

    query := `INSERT INTO friends (user_a_id, user_b_id) VALUES ($1, $2)`
	_, err = c.db.Exec(query, recipientId, senderId)
	if err != nil {
		return fmt.Errorf("failed to accept friend request: %w", err)
	}
	return nil
}

func (c *PostgresFriendsRepository) DeleteRequest(senderId uuid.UUID, recipientId uuid.UUID) error{
    query := `DELETE FROM friends_requests WHERE recipient_id= $1 AND sender_id= $2`
    _, err := c.db.Exec(query, recipientId, senderId)
	if err != nil {
		return fmt.Errorf("failed to delete friend request: %w", err)
	}
    return nil
}

func (c PostgresFriendsRepository) RejectFriendRequest(senderId uuid.UUID, recipientId uuid.UUID) error {
    err := c.DeleteRequest(senderId, recipientId)
    if err != nil {
        return fmt.Errorf("failed to reject friend request: %w", err)
    }
    return nil
}

func (c PostgresFriendsRepository) GetFriendRequestsSent(senderId uuid.UUID) ([]um.UserResponse, error) {
	friends := []um.UserRecord{}
	query := `
        SELECT 
            us.* 
        FROM friends_requests fr
        JOIN users us ON us.id = fr.recipient_id
        WHERE fr.sender_id = $1
    `
	err := c.db.Select(&friends, query, senderId)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get friend requests sent: %w", err)
	}

    res := []um.UserResponse{}
    for _, friend := range friends {
        res = append(res, *utils.MapUserRecordToUserResponse(&friend))
    }
	return res, nil
}

func (c PostgresFriendsRepository) GetFriendRequestsReceived(receiverId uuid.UUID) ([]um.UserResponse, error) {
	friends := []um.UserRecord{}
	query := `
        SELECT 
            us.* 
        FROM friends_requests fr
        JOIN users us ON us.id = fr.sender_id
        WHERE fr.recipient_id= $1
    `
	err := c.db.Select(&friends, query, receiverId)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get friend requests sent: %w", err)
	}

    res := []um.UserResponse{}
    for _, friend := range friends {
        res = append(res, *utils.MapUserRecordToUserResponse(&friend))
    }
	return res, nil
}

func (c PostgresFriendsRepository) CheckIfFriendRequestExists(senderId uuid.UUID, recipientId uuid.UUID) bool {
	query := `SELECT EXISTS (SELECT 1 FROM friends_requests WHERE recipient_id= $1 AND sender_id= $2)`
	var exists1 bool
	err := c.db.Get(&exists1, query, recipientId, senderId)
	if err != nil {
		return false
	}

	return exists1 
}

func (c PostgresFriendsRepository) CheckIfFriendShipExists(userA uuid.UUID, userB uuid.UUID) bool {
	query := `SELECT EXISTS (SELECT 1 FROM friends WHERE (user_a_id = $1 AND user_b_id= $2) OR (user_a_id = $2 AND user_b_id= $1))`
	var exists1 bool
	err := c.db.Get(&exists1, query, userA, userB)
	if err != nil {
		return false
	}
	return exists1 
}


func (c PostgresFriendsRepository) DeleteFriendship(userA uuid.UUID, userB uuid.UUID) error{
    query := `DELETE FROM friends WHERE (user_a_id = $1 AND user_b_id= $2) OR (user_a_id = $2 AND user_b_id= $1)`
    _, err := c.db.Exec(query, userA, userB)
    if err != nil {
        return fmt.Errorf("failed to delete friendship: %w", err)
    }
    return nil
}
