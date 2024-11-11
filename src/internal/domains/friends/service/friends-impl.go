package service

import (
    "github.com/google/uuid"
    "github.com/betterreads/internal/domains/friends/repository"
    "github.com/betterreads/internal/domains/friends/models"
    users "github.com/betterreads/internal/domains/users/service"
)

type FriendsServiceImpl struct{
    fr repository.FriendsRepository
    us users.UsersService
}

func NewFriendsServiceImpl(fr repository.FriendsRepository, us users.UsersService) FriendsService{
    return &FriendsServiceImpl{fr: fr, us: us}
}



func (fs *FriendsServiceImpl) GetFriends(userID uuid.UUID) ([]models.FriendOfUser, error){
    if !fs.us.CheckUserExists(userID) {
        return nil, users.ErrUserNotFound
    }


    friends , err := fs.fr.GetFriends(userID)
    if err != nil {
        return nil, err
    }
    return friends, nil
}

func (fs *FriendsServiceImpl) AddFriend(userID uuid.UUID, friendID uuid.UUID) error{
    if userID == friendID {
        return ErrSameUser
    }
    if !fs.us.CheckUserExists(friendID) {
        return ErrUserFriendNotFound
    }

    if !fs.us.CheckUserExists(userID) {
        return users.ErrUserNotFound
    }

    if fs.fr.CheckIfFriendRequestExists(userID, friendID) {
        return ErrFriendRequestExists
    }

    if fs.fr.CheckIfFriendShipExists(friendID, userID) {
        return ErrAlreadyFriends
    }

    err := fs.fr.AddFriend(userID, friendID)
    if err != nil {
        return err
    }
    return nil
}

func (fs *FriendsServiceImpl) AcceptFriendRequest(userID uuid.UUID, friendID uuid.UUID) error{
    if !fs.fr.CheckIfFriendRequestExists(userID, friendID) {
        return ErrRequestNotFound
    }

    err := fs.fr.AcceptFriendRequest(userID, friendID)
    if err != nil {
        return err
    }
    return nil
}

func (fs *FriendsServiceImpl) RejectFriendRequest(userID uuid.UUID, friendID uuid.UUID) error{
    if !fs.fr.CheckIfFriendRequestExists(userID, friendID) {
        return ErrRequestNotFound
    }


    err := fs.fr.RejectFriendRequest(userID, friendID)
    if err != nil {
        return err
    }
    return nil
}

func (fs *FriendsServiceImpl) GetFriendRequestsSent(userID uuid.UUID) ([]models.FriendOfUser, error){
    if !fs.us.CheckUserExists(userID) {
        return nil, users.ErrUserNotFound
    }

    friendRequestsSent, err := fs.fr.GetFriendRequestsSent(userID)
    if err != nil {
        return nil, err
    }
    return friendRequestsSent, nil
}

func (fs *FriendsServiceImpl) GetFriendRequestsReceived(userID uuid.UUID) ([]models.FriendOfUser, error){
    if !fs.us.CheckUserExists(userID) {
        return nil, users.ErrUserNotFound
    }


    friendRequestsReceived, err := fs.fr.GetFriendRequestsReceived(userID)
    if err != nil {
        return nil, err
    }
    return friendRequestsReceived, nil
}

