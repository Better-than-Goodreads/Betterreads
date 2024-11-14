package service

import (
	"github.com/betterreads/internal/domains/friends/repository"
	"github.com/betterreads/internal/domains/users/models"
	users "github.com/betterreads/internal/domains/users/service"
	"github.com/google/uuid"
)

type FriendsServiceImpl struct {
	fr repository.FriendsRepository
	us users.UsersService
}

func NewFriendsServiceImpl(fr repository.FriendsRepository, us users.UsersService) FriendsService {
	return &FriendsServiceImpl{fr: fr, us: us}
}

func (fs *FriendsServiceImpl) GetFriends(userID uuid.UUID) ([]models.UserResponse, error) {
	if !fs.us.CheckUserExists(userID) {
		return nil, users.ErrUserNotFound
	}

	friends, err := fs.fr.GetFriends(userID)
	if err != nil {
		return nil, err
	}
	return friends, nil
}

func (fs *FriendsServiceImpl) AddFriend(senderId uuid.UUID, recipientId uuid.UUID) error {
	if senderId == recipientId {
		return ErrSameUser
	}
	if !fs.us.CheckUserExists(senderId) {
		return ErrUserFriendNotFound
	}

	if !fs.us.CheckUserExists(recipientId) {
		return users.ErrUserNotFound
	}

	if fs.fr.CheckIfFriendRequestExists(senderId, recipientId) {
		return ErrFriendRequestExists
	}

	if fs.fr.CheckIfFriendShipExists(senderId, recipientId) {
		return ErrAlreadyFriends
	}

	err := fs.fr.AddFriend(senderId, recipientId)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FriendsServiceImpl) AcceptFriendRequest(recipientId uuid.UUID, senderId uuid.UUID) error {
	if !fs.fr.CheckIfFriendRequestExists(senderId, recipientId) {
		return ErrRequestNotFound
	}

	err := fs.fr.AcceptFriendRequest(recipientId, senderId)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FriendsServiceImpl) RejectFriendRequest(recipientId uuid.UUID, senderId uuid.UUID) error {
	if !fs.fr.CheckIfFriendRequestExists(senderId, recipientId) {
		return ErrRequestNotFound
	}

	err := fs.fr.RejectFriendRequest(recipientId, senderId)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FriendsServiceImpl) GetFriendRequestsSent(userID uuid.UUID) ([]models.UserResponse, error) {
	if !fs.us.CheckUserExists(userID) {
		return nil, users.ErrUserNotFound
	}

	friendRequestsSent, err := fs.fr.GetFriendRequestsSent(userID)
	if err != nil {
		return nil, err
	}
	return friendRequestsSent, nil
}

func (fs *FriendsServiceImpl) GetFriendRequestsReceived(userID uuid.UUID) ([]models.UserResponse, error) {
	if !fs.us.CheckUserExists(userID) {
		return nil, users.ErrUserNotFound
	}

	friendRequestsReceived, err := fs.fr.GetFriendRequestsReceived(userID)
	if err != nil {
		return nil, err
	}
	return friendRequestsReceived, nil
}

func (fs *FriendsServiceImpl) DeleteFriend(userID uuid.UUID, friendID uuid.UUID) error {
	if userID == friendID {
		return ErrSameUser
	}
	if !fs.us.CheckUserExists(friendID) {
		return ErrUserFriendNotFound
	}

	if !fs.us.CheckUserExists(userID) {
		return users.ErrUserNotFound
	}

	if !fs.fr.CheckIfFriendShipExists(friendID, userID) {
		return ErrFriendShipNotFound
	}

	if err := fs.fr.DeleteFriendship(userID, friendID); err != nil {
		return err
	}

	return nil
}
