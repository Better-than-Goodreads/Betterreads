package repository

import (
	_ "database/sql"
	_ "github.com/lib/pq"
	"fmt"

	"github.com/betterreads/internal/domains/users/models"
	"github.com/jmoiron/sqlx"
)


type PostgresUserRepository struct {
	c *sqlx.DB
}

func NewPostgresUserRepository(c *sqlx.DB) (*PostgresUserRepository, error) {
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
			about_me TEXT
		);
		
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	schemaRegistry := `
		CREATE TABLE IF NOT EXISTS registry (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			username VARCHAR(255) NOT NULL UNIQUE,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
		
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	if _, err := c.Exec(schemaUsers); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	if _, err := c.Exec(schemaRegistry); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PostgresUserRepository{c} , nil
}

func (r *PostgresUserRepository) CreateStageUser(user *models.UserStageRequest) (*models.UserStageRecord, error) { 
	// err := r.checkUserExists(user)
	
	var userRecord *models.UserStageRecord
	query := `INSERT INTO registry (email, username, password, first_name, last_name)
			  VALUES (:username, :email, :password, :first_name, :last_name) 
			  RETURNING id, email, username, first_name, last_name;`

	rows , err := r.c.NamedQuery(query, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(userRecord); err != nil {
			return nil, fmt.Errorf("error scanning user data: %w", err)
		}
	} else {
		return nil, fmt.Errorf("error: no user created")
	}

	return userRecord, nil
}

func (r *PostgresUserRepository) GetStageUser(uuid string) (*models.UserStageRecord, error) {
	var user *models.UserStageRecord
	query := `SELECT * FROM registry WHERE id = $1;`
	if err := r.c.Select(user, query, uuid); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}


func (r *PostgresUserRepository) JoinAndCreateUser(userAddtional *models.UserAdditionalRequest) (*models.UserRecord, error) {
	user , err := r.GetStageUser(userAddtional.Id)	
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	userRecord := &models.UserRecord{
		Email: user.Email,
		Password: user.Password,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Username: user.Username,
		Location: userAddtional.Location,
		Age: userAddtional.Age,
		Gender: userAddtional.Gender,
		AboutMe: userAddtional.AboutMe,
	}

	if err := r.createUser(userRecord) ; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := r.deleteStageUser(userAddtional.Id); err != nil {
		return nil, fmt.Errorf("failed to delete stage user: %w", err)
	}

	return userRecord, nil
}

func (r *PostgresUserRepository) createUser(user *models.UserRecord) error {
	userRecord := &models.UserRecord{}

	query := `INSERT INTO users (email, username, password, first_name, last_name,username, location, age, gender, about_me)
	          VALUES (:email, :username, :password, :first_name, :last_name, :username, :location, :age, :gender, :about_me)
			  RETURNING email, username, first_name, last_name, username, location, age, gender, about_me;`


	if err := r.c.QueryRowx(query, user).StructScan(userRecord); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}	
	return nil
}

func (r *PostgresUserRepository) deleteStageUser(id string) error {
	query := `DELETE FROM registry WHERE id = $1;`
	if _, err := r.c.Exec(query, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *PostgresUserRepository) GetUser(id string) (*models.UserRecord, error) {
	var user *models.UserRecord
	query := `SELECT * FROM users WHERE id = $1;`
	if err := r.c.Select(user, query, id); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) GetUsers() ([]*models.UserRecord, error) {
	var users []*models.UserRecord
	query := `SELECT * FROM users;`
	if err := r.c.Select(users, query); err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

func (r *PostgresUserRepository) GetUserByUsername(username string) (*models.UserRecord, error) {
	var user *models.UserRecord
	query := `SELECT * FROM users WHERE username = $1;`
	if err := r.c.Select(user, query, username); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (*models.UserRecord, error) {
	var user *models.UserRecord
	query := `SELECT * FROM users WHERE email = $1;`
	if err := r.c.Select(user, query, email); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// func (r *PostgresUserRepository) checkUserExists(user *models.UserStageRequest) error {
// 	var userRecord *models.UserRecord
// 	query := `SELECT * FROM users WHERE email = $1 OR username = $2;`
// 	if err := r.c.Select(userRecord, query, user.Email, user.Username); err != nil {
// 		return  fmt.Errorf("failed to get user: %w", err)
// 	}

// 	var userStage *models.UserStageRecord 
// 	query = `SELECT * FROM registry WHERE email = $1 OR username = $2;`
// 	if err := r.c.Select(userStage, query, user.Email, user.Username); err != nil {
// 		return  fmt.Errorf("failed to get user: %w", err)
// 	}
	
// 	return nil
// }

