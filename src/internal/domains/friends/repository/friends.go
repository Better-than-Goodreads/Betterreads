package repository

import (
	um "github.com/betterreads/internal/domains/users/models"
	"github.com/google/uuid"
)

type FriendsRepository interface {
	GetFriends(userID uuid.UUID) ([]um.UserResponse, error)
	AddFriend(sender uuid.UUID, recipient uuid.UUID) error
	AcceptFriendRequest(recipient uuid.UUID, sender uuid.UUID) error
	RejectFriendRequest(recipient uuid.UUID, sender uuid.UUID) error
	GetFriendRequestsSent(sender uuid.UUID) ([]um.UserResponse, error)
	GetFriendRequestsReceived(recipient uuid.UUID) ([]um.UserResponse, error)
	CheckIfFriendRequestExists(sender uuid.UUID, recipient uuid.UUID) bool
	CheckIfFriendShipExists(userA uuid.UUID, userB uuid.UUID) bool
    DeleteFriendship(userA uuid.UUID, userB uuid.UUID) error
}
