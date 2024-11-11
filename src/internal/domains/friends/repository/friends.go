package repository

import (
	"github.com/betterreads/internal/domains/friends/models"
	"github.com/google/uuid"
)

type FriendsRepository interface {
	GetFriends(userID uuid.UUID) ([]models.FriendOfUser, error)
	AddFriend(userID uuid.UUID, friendID uuid.UUID) error
	AcceptFriendRequest(userID uuid.UUID, friendID uuid.UUID) error
	RejectFriendRequest(userID uuid.UUID, friendID uuid.UUID) error
	GetFriendRequestsSent(userID uuid.UUID) ([]models.FriendOfUser, error)
	GetFriendRequestsReceived(userID uuid.UUID) ([]models.FriendOfUser, error)
	CheckIfFriendRequestExists(userID uuid.UUID, friendID uuid.UUID) bool
	CheckIfFriendShipExists(userID uuid.UUID, friendID uuid.UUID) bool
}
