package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/betterreads/internal/domains/communities/model"
	"github.com/google/uuid"
	userModel "github.com/betterreads/internal/domains/users/models"
)

type PostgresCommunitiesRepository struct {
	db *sqlx.DB
}

func NewPostgresCommunitiesRepository(db *sqlx.DB) (CommunitiesDatabase, error) {
	enableUUIDExtension := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	if _, err := db.Exec(enableUUIDExtension); err != nil {
		return nil, fmt.Errorf("failed to enable uuid extension: %w", err)
	}

	schemaCommunities := `
		CREATE TABLE IF NOT EXISTS communities (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			owner_id UUID NOT NULL,
			FOREIGN KEY (owner_id) REFERENCES users(id)
		);`

	if _, err := db.Exec(schemaCommunities); err != nil {
		return nil, fmt.Errorf("failed to create communities table: %w", err)
	}

	schemaCommunitiesUsers := `
		CREATE TABLE IF NOT EXISTS communities_users (
			user_id UUID NOT NULL,
			community_id UUID NOT NULL,
			PRIMARY KEY (user_id, community_id),	
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (community_id) REFERENCES communities(id)
		);`
	
	if _, err := db.Exec(schemaCommunitiesUsers); err != nil {
		return nil, fmt.Errorf("failed to create communities_users table: %w", err)
	}

	return &PostgresCommunitiesRepository{db: db}, nil
}

func (db *PostgresCommunitiesRepository) CreateCommunity(community model.NewCommunityRequest, userId uuid.UUID) (*model.CommunityResponse, error) {
	query := `INSERT INTO communities (name, description, owner_id) VALUES ($1, $2, $3) RETURNING id`
	
	var id uuid.UUID
	err := db.db.QueryRow(query, community.Name, community.Description, userId).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create community: %w", err)
	}
	
	communityResponse := model.CommunityResponse{
		ID: id,
		Name: community.Name,
		Description: community.Description,
		OwnerID: userId,
	}
	fmt.Println(communityResponse)
	return &communityResponse, nil
}

func (db *PostgresCommunitiesRepository) GetCommunities() ([]*model.CommunityResponse, error) {
	query := `SELECT id, name, description, owner_id FROM communities`
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get communities: %w", err)
	}
	defer rows.Close()

	communities := []*model.CommunityResponse{}
	for rows.Next() {
		community := &model.CommunityResponse{}

		err := rows.Scan(&community.ID, &community.Name, &community.Description, &community.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan community: %w", err)
		}
		communities = append(communities, community)
	}

	return communities, nil
}

func (db *PostgresCommunitiesRepository) JoinCommunity(communityId uuid.UUID, userId uuid.UUID) error {
	query := `INSERT INTO communities_users (user_id, community_id) VALUES ($1, $2)`
	_, err := db.db.Exec(query, userId, communityId)
	if err != nil {
		return fmt.Errorf("failed to join community: %w", err)
	}

	return nil
}
func (db *PostgresCommunitiesRepository) CheckIfUserIsInCommunity(communityId uuid.UUID, userId uuid.UUID) bool {
	query := `SELECT EXISTS(SELECT 1 FROM communities_users WHERE user_id=$1 AND community_id=$2)`

	var exists bool
	err := db.db.QueryRow(query, userId, communityId).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}


func (db *PostgresCommunitiesRepository) GetCommunityUsers(communityId uuid.UUID) ([]*userModel.UserStageResponse, error) {
	query := `SELECT u.email, u.username, u.first_name, u.last_name, u.is_author, u.id FROM users u 
			  JOIN communities_users cu ON u.id = cu.user_id 
			  WHERE cu.community_id = $1`
	rows, err := db.db.Query(query, communityId)
	if err != nil {
		return nil, fmt.Errorf("failed to get community users: %w", err)
	}
	defer rows.Close()

	users := []*userModel.UserStageResponse{}
	for rows.Next() {
		user := &userModel.UserStageResponse{}

		err := rows.Scan(&user.Email, &user.Username, &user.First_name, &user.Last_name, &user.IsAuthor, &user.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}