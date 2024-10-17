package repository

import (
	"github.com/betterreads/internal/domains/users/models"
	"github.com/betterreads/internal/domains/users/utils"
)

type MemoryDatabase struct {
	users map[int]models.UserRecord
	currId	int
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		users: make(map[int]models.UserRecord),
	}
}

func (m *MemoryDatabase) CreateUser(user models.UserRequest) (models.UserRecord, error) {
	id := m.currId + 1

	userRecord := utils.MapUserRequestToUserRecord(user,id)

	m.users[id] = userRecord
	return userRecord, nil
}

func (m *MemoryDatabase) GetUser(id int) (models.UserRecord, error) {
	user, ok := m.users[id]
	if !ok {
		return models.UserRecord{}, ErrUserNotFound 
	}
	return user, nil
}


func (m *MemoryDatabase) GetUsers() ([]models.UserRecord, error) {
	usersArr := make([]models.UserRecord, len(m.users))
	for _, v := range m.users{
		usersArr = append(usersArr, v)
	}
	return usersArr, nil
}
