package repository

import (
	"github.com/betterreads/internal/domains/users/models"
	"github.com/betterreads/internal/domains/users/utils"
)

type MemoryDatabase struct {
	users  map[int]*models.UserRecord
	currId int
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		users: make(map[int]*models.UserRecord),
	}
}

func (m *MemoryDatabase) CreateUser(user *models.UserRequest) (*models.UserRecord, error) {
	id := m.currId + 1

	userRecord := utils.MapUserRequestToUserRecord(user, id)

	m.users[id] = userRecord
	m.currId = id
	return userRecord, nil
}

func (m *MemoryDatabase) GetUser(id int) (*models.UserRecord, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (m *MemoryDatabase) GetUsers() ([]*models.UserRecord, error) {
	usersArr := make([]*models.UserRecord, 0, len(m.users))
	for _, user := range m.users {
		usersArr = append(usersArr, user)
	}
	return usersArr, nil
}

func (m *MemoryDatabase) GetUserByUsername(username string) (*models.UserRecord, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (m *MemoryDatabase) GetUserByEmail(email string) (*models.UserRecord, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}
