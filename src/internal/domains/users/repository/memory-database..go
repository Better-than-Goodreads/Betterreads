package repository

import (
	"github.com/betterreads/internal/domains/users/models"
	"github.com/betterreads/internal/domains/users/utils"
    "github.com/google/uuid"
)

type MemoryDatabase struct {
	users  map[int]*models.UserRecord
    registeringUsers map[string]*models.UserStageRecord
	currId int
}

func NewMemoryDatabase() UsersDatabase{
    db := new(MemoryDatabase)
    db.users = make(map[int]*models.UserRecord)
    db.registeringUsers = make(map[string]*models.UserStageRecord)
    db.currId = 0
    return db
}

func (m *MemoryDatabase) createUser(user *models.UserRequest) (*models.UserRecord, error) {
 	id := m.currId + 1

 	userRecord := utils.MapUserRequestToUserRecord(user, id)

 	m.users[id] = userRecord
 	m.currId = id
 	return userRecord, nil
}

func (m *MemoryDatabase) CreateStageUser(user *models.UserStageRequest) (*models.UserStageRecord, error) {
    if err := m.checkUserExist(user.Username, user.Email); err != nil {
        return nil, err
    }

    token := uuid.New().String()
    userRecord := utils.MapUserStageRequestToUserStageRecord(user, token)

    m.registeringUsers[token] = userRecord 
    return userRecord, nil
}

func (m *MemoryDatabase) JoinAndCreateUser (userAdditional *models.UserAdditionalRequest) (*models.UserRecord, error) {
    user, ok := m.registeringUsers[userAdditional.Token]
    if !ok {
        return nil, ErrUserNotFound
    }
    UserRequest := &models.UserRequest{
        Email: user.Email,
        Password: user.Password,
        FirstName: user.FirstName,
        LastName: user.LastName,
        Username: user.Username,
        Location: userAdditional.Location,
        Gender: userAdditional.Gender,
        AboutMe: userAdditional.AboutMe,
    }

    userRecord, err := m.createUser(UserRequest)
    return userRecord, err
}

func (m *MemoryDatabase) deleteStageUser(id int) error {
    _, ok := m.users[id]
    if !ok {
        return ErrUserNotFound
    }
    delete(m.users, id)
    return nil
}

func (m *MemoryDatabase) checkUserExist(username string, email string) error {
    for _, user := range m.users {
        if user.Username == username {
            return ErrUsernameAlreadyTaken
        }
        if user.Email == email {
            return ErrEmailAlreadyTaken
        }
    }

    for _, user := range m.registeringUsers {
        if user.Username == username {
            return ErrUsernameAlreadyTaken
            }
        if user.Email == email {
            return ErrEmailAlreadyTaken
        }
    }

    return nil 
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

func (m *MemoryDatabase) GetStageUser(token string) (*models.UserStageRecord, error) {
    user, ok := m.registeringUsers[token]
    if !ok {
        return nil, ErrUserNotFound
    }
    return user, nil
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


func (m *MemoryDatabase) CreateStagingUser(user *models.UserStageRecord) (*models.UserStageRecord, error) {
    m.registeringUsers[user.Username] = user
    return user, nil
}
