package service

import ( 
    "github.com/google/uuid"
    "github.com/betterreads/internal/domains/friends/models"
    "errors"
)

var (
    ErrUserFriendNotFound = errors.New("user friend not found")
    ErrFriendRequestExists = errors.New("friend request already exists")
    ErrAlreadyFriends = errors.New("users are already friends")
    ErrRequestNotFound = errors.New("friend request not found")
    ErrSameUser = errors.New("cannot add yourself as a friend")
)

type FriendsService interface {
    GetFriends(userID uuid.UUID) ([]models.FriendOfUser, error)
    AddFriend(userID uuid.UUID, friendID uuid.UUID) error
    AcceptFriendRequest(userID uuid.UUID, friendID uuid.UUID) error
    RejectFriendRequest(userID uuid.UUID, friendID uuid.UUID) error
    GetFriendRequestsSent(userID uuid.UUID) ([]models.FriendOfUser, error)
    GetFriendRequestsReceived(userID uuid.UUID) ([]models.FriendOfUser, error)
}
