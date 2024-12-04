package service

import (
	"errors"
	"github.com/betterreads/internal/domains/users/models"
	"github.com/google/uuid"
)

var (
	ErrUserFriendNotFound  = errors.New("user friend not found")
	ErrFriendRequestExists = errors.New("friend request already exists")
    ErrFriendShipNotFound  = errors.New("friendship not found")
	ErrAlreadyFriends      = errors.New("users are already friends")
	ErrRequestNotFound     = errors.New("friend request not found")
	ErrSameUser            = errors.New("cannot add yourself as a friend")
)

type FriendsService interface {
	GetFriends(userID uuid.UUID) ([]models.UserResponse, error)
	AddFriend(senderId uuid.UUID, recipientId uuid.UUID) error
	AcceptFriendRequest(recipientId uuid.UUID, senderId uuid.UUID) error
	RejectFriendRequest(recipientId uuid.UUID, senderId uuid.UUID) error
	GetFriendRequestsSent(userID uuid.UUID) ([]models.UserResponse, error)
	GetFriendRequestsReceived(userID uuid.UUID) ([]models.UserResponse, error)
    DeleteFriend(userA uuid.UUID, userB uuid.UUID) error
}
