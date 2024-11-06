package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
    "github.com/google/uuid"
	"github.com/betterreads/internal/domains/users/models"
	"github.com/jmoiron/sqlx"
)



type PostgresUserRepository struct {
	c *sqlx.DB
}

func NewPostgresUserRepository(c *sqlx.DB) (UsersDatabase, error) {
	enableUUIDExtension := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	if _, err := c.Exec(enableUUIDExtension); err != nil {
		return nil, fmt.Errorf("failed to enable uuid extension: %w", err)
	}

	schemaUsers := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(255) NOT NULL UNIQUE,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			location VARCHAR(255) NULL,
			age INTEGER,
			gender VARCHAR(255),
			about_me TEXT,
            is_author BOOLEAN DEFAULT FALSE
		);
		
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	schemaRegistry := `
		CREATE TABLE IF NOT EXISTS registry (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email VARCHAR(255) NOT NULL UNIQUE,
			username VARCHAR(255) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
            is_author BOOLEAN DEFAULT FALSE
		);
		
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	schemaPictures := `
		CREATE TABLE IF NOT EXISTS pictures_users (
			user_id UUID,
			picture BYTEA,
			FOREIGN KEY (user_id) REFERENCES users(id),
			PRIMARY KEY (user_id)
		);
		`

	if _, err := c.Exec(schemaUsers); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	if _, err := c.Exec(schemaRegistry); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

    if _, err := c.Exec(schemaPictures); err != nil {
        return nil, fmt.Errorf("failed to create table: %w", err)
    }

    

	return &PostgresUserRepository{c}, nil
}

func (r *PostgresUserRepository) CreateStageUser(user *models.UserStageRequest) (*models.UserStageRecord, error) {
    userRecord := &models.UserStageRecord{}
	query :=`INSERT INTO registry (email, username, password, first_name, last_name, 
            is_author)
            VALUES ($1, $2, $3, $4, $5, $6)
            RETURNING id, email, username, first_name, last_name, is_author;`

    args := []interface{}{user.Email, user.Username, user.Password, user.FirstName, user.LastName, user.IsAuthor}
	
    err := r.c.Get(userRecord, query, args...)
    if err != nil {
        return nil,  err
    }
    
	return userRecord, nil
}

func (r *PostgresUserRepository) JoinAndCreateUser(userAdditional *models.UserAdditionalRequest, id uuid.UUID) (*models.UserRecord, error) {
	user, err := r.GetStageUser(id)
	if err != nil {
		return nil, err
	}
    userRecord, err := r.createUser(user, userAdditional)
    if err != nil {
		return nil, err
	}
    

	if err := r.deleteStageUser(id); err != nil {
		return nil, err
	}

	return userRecord, nil
}

func (r *PostgresUserRepository) createUser(user *models.UserStageRecord, userAdditional *models.UserAdditionalRequest) (*models.UserRecord, error) {
    userRecord := &models.UserRecord{}
	query := `INSERT INTO users (email, password, first_name, last_name, username, 
                    location, gender, about_me, age, is_author)
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                    RETURNING id, email, password, first_name, last_name, username, location, gender,about_me, age, is_author;`
    
    args := []interface{}{user.Email, user.Password, user.FirstName, user.LastName, user.Username, userAdditional.Location,
					userAdditional.Gender, userAdditional.AboutMe, userAdditional.Age, user.IsAuthor} 
    
    err := r.c.Get(userRecord, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return userRecord, nil
}

func (r *PostgresUserRepository) deleteStageUser(id uuid.UUID) error {
	query := `DELETE FROM registry WHERE id = $1;`
	if _, err := r.c.Exec(query, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *PostgresUserRepository) GetUser(id uuid.UUID) (*models.UserRecord, error) {
    user := &models.UserRecord{}
	query := `SELECT * FROM users WHERE id = $1;`
	if err := r.c.Get(user, query, id); err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) GetStageUser(uuid uuid.UUID) (*models.UserStageRecord, error) {
    user := &models.UserStageRecord{}
	query := `SELECT * FROM registry WHERE id = $1;`
	if err := r.c.Get(user, query, uuid); err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserStageNotFound
        }
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) GetUsers() ([]*models.UserRecord, error) {
    users := []*models.UserRecord{}
	query := `SELECT * FROM users;`
	if err := r.c.Select(&users, query); err != nil {
        if err != sql.ErrNoRows {
            return nil, fmt.Errorf("failed to get users: %w", err)
        }
	}
	return users, nil
}

func (r *PostgresUserRepository) GetUserByUsername(username string) (*models.UserRecord, error) {
    user := &models.UserRecord{}
	query := `SELECT * FROM users WHERE username = $1;`
	if err := r.c.Get(user, query, username); err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (*models.UserRecord, error) {
    user := &models.UserRecord{}
	query := `SELECT * FROM users WHERE email = $1;`
	if err := r.c.Get(user, query, email); err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) GetUserPicture(id uuid.UUID) ([]byte, error) {
	var picture []byte
	query := `SELECT picture FROM pictures_users WHERE user_id= $1;`
	if err := r.c.Get(&picture, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return picture, nil
}

func (r *PostgresUserRepository) SaveUserPicture(id uuid.UUID, picture []byte) error {
    exists := false
    query := `SELECT EXISTS(SELECT 1 FROM pictures_users WHERE user_id = $1);`
    if err := r.c.Get(&exists, query, id); err != nil {
        if err == sql.ErrNoRows {
            return ErrUserNotFound
        }
        return fmt.Errorf("failed to check if user picture exists: %w", err)
    }

    if exists {
        query = `UPDATE pictures_users SET picture = $2 WHERE user_id = $1;`

    } else {
        query = `INSERT INTO pictures_users (user_id, picture) VALUES ($1, $2);`
    }

    if _, err := r.c.Exec(query, id, picture); err != nil {
        return fmt.Errorf("failed to save user picture: %w", err)
    }
    return nil
}

func (r *PostgresUserRepository) CheckUserExistsForRegister(user *models.UserStageRequest) error {
    result_registry_email, result_registry_username , err:= r.checkUserExistsInTable(user , "registry")
    if err != nil {
        return  fmt.Errorf("failed to check user exists in user table: %w", err)
    }
    result_users_email, result_users_username, err :=  r.checkUserExistsInTable(user , "users")
    if err != nil {
        return  fmt.Errorf("failed to check user exists in user table: %w", err)
    }


    if result_registry_email || result_users_email {
        return ErrEmailAlreadyTaken
    }

    if result_registry_username || result_users_username {
        return ErrUsernameAlreadyTaken
    }


    return nil 
}
    
func (r *PostgresUserRepository) checkUserExistsInTable(user *models.UserStageRequest, table string) (bool, bool, error) {
    result_email := false
    result_username := false

    // Check if the email exists
    query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE email = $1);`, table)
    err := r.c.Get(&result_email, query, user.Email)
    if err != nil {
        return false, false, fmt.Errorf("failed to check user exists for email: %w", err)
    }

    // Check if the username exists
    query = fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE username = $1);`, table)
    err = r.c.Get(&result_username, query, user.Username)
    if err != nil {
        return false, false, fmt.Errorf("failed to check user exists for username: %w", err)
    }

    return result_email, result_username, nil
}

func (r *PostgresUserRepository) CheckUserExists(id uuid.UUID) bool {
    exists := false
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);`
    if err := r.c.Get(&exists, query, id); err != nil {
        return false
    }
    return exists
}


func (r *PostgresUserRepository) SearchUsers(username string, isAuthor bool) ([]*models.UserRecord, error) {
	users := []*models.UserRecord{}
	query := `SELECT * FROM users WHERE LOWER(username) LIKE LOWER('%'||$1||'%') AND is_author = $2;`
	if err := r.c.Select(&users, query, username, isAuthor); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to search users: %w", err)
		}
	}
	return users, nil
}