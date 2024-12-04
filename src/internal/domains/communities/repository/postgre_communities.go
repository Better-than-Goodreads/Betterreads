package repository

import (
	"database/sql"
	"fmt"

	"github.com/betterreads/internal/domains/communities/model"
	userModel "github.com/betterreads/internal/domains/users/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

	schemaCommunitiesPictures := `
		CREATE TABLE IF NOT EXISTS communities_pictures (
			community_id UUID NOT NULL,
			picture BYTEA NOT NULL,
			PRIMARY KEY (community_id),
			FOREIGN KEY (community_id) REFERENCES communities(id)
		);
	`

	if _, err := db.Exec(schemaCommunitiesPictures); err != nil {
		return nil, fmt.Errorf("failed to create communities_pictures table: %w", err)
	}

	schemaCommunitiesPosts := `
		CREATE TABLE IF NOT EXISTS communities_posts (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			community_id UUID NOT NULL,
			user_id UUID NOT NULL,
			content TEXT NOT NULL,
			title VARCHAR(255) NOT NULL,
			date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (community_id) REFERENCES communities(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`

	if _, err := db.Exec(schemaCommunitiesPosts); err != nil {
		return nil, fmt.Errorf("failed to create communities_posts table: %w", err)
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

	query = `INSERT INTO communities_pictures (community_id, picture) VALUES ($1, $2)`
	_, err = db.db.Exec(query, id, community.Picture)
	if err != nil {
		return nil, fmt.Errorf("failed to create community picture: %w", err)
	}

	communityResponse := model.CommunityResponse{
		ID:          id,
		Name:        community.Name,
		Description: community.Description,
		OwnerID:     userId,
		Joined:      true,
	}

	JoinCommunityErr := db.JoinCommunity(id, userId)
	if JoinCommunityErr != nil {
		return nil, fmt.Errorf("failed to join user to community: %w", JoinCommunityErr)
	}

	return &communityResponse, nil
}

func (db *PostgresCommunitiesRepository) GetCommunities(userId uuid.UUID) ([]*model.CommunityResponse, error) {
	query := `SELECT 
    c.id AS id, 
    c.name AS name, 
    c.description AS desc, 
    c.owner_id AS owner,
    CASE 
        WHEN cu.user_id IS NOT NULL THEN true 
        ELSE false 
    END AS joined
	FROM 
    communities c
	LEFT JOIN 
    communities_users cu 
    ON c.id = cu.community_id AND cu.user_id = $1`
	rows, err := db.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get communities: %w", err)
	}
	defer rows.Close()

	communities := []*model.CommunityResponse{}
	for rows.Next() {
		community := &model.CommunityResponse{}

		err := rows.Scan(&community.ID, &community.Name, &community.Description, &community.OwnerID, &community.Joined)
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

func (db *PostgresCommunitiesRepository) GetCommunityPicture(communityId uuid.UUID) ([]byte, error) {
	query := `SELECT picture FROM communities_pictures WHERE community_id = $1`
	var picture []byte
	err := db.db.QueryRow(query, communityId).Scan(&picture)
	if err != nil {
		if err == sql.ErrNoRows {
			return picture, nil
		}
		return nil, fmt.Errorf("failed to get community picture: %w", err)
	}

	return picture, nil
}

func (db *PostgresCommunitiesRepository) SearchCommunities(search string, curr_user uuid.UUID) ([]*model.CommunityResponse, error) {
	query := `SELECT 
    c.id, 
    c.name, 
    c.description,
    c.owner_id,  -- Add a comma here
    CASE 
        WHEN cu.user_id IS NOT NULL THEN true 
        ELSE false 
    END AS joined
    FROM communities c
    LEFT JOIN communities_users cu 
        ON c.id = cu.community_id AND cu.user_id = $2
    WHERE c.name ILIKE '%' || $1 || '%'`

	var communities []*model.CommunityResponse
	err := db.db.Select(&communities, query, search, curr_user)
	if err != nil {
		return nil, fmt.Errorf("failed to search communities: %w", err)
	}

	return communities, nil
}

func (db *PostgresCommunitiesRepository) GetCommunityById(id uuid.UUID, userId uuid.UUID) (*model.CommunityResponse, error) {
	query := `SELECT 
    c.id, 
    c.name, 
    c.description,
    c.owner_id,  -- Add a comma here
    CASE 
        WHEN cu.user_id IS NOT NULL THEN true 
        ELSE false 
    END AS joined
    FROM communities c
    LEFT JOIN communities_users cu 
        ON c.id = cu.community_id and cu.user_id = $1
    WHERE c.id = $2 `

	var community model.CommunityResponse
	err := db.db.Get(&community, query, userId, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCommunityNotFound
		}
		return nil, fmt.Errorf("failed to get community: %w", err)
	}
	return &community, nil
}

func (db *PostgresCommunitiesRepository) GetCommunityPosts(communityId uuid.UUID) ([]*model.CommunityPostResponse, error) {
	query := `SELECT 
	cp.id, 
	cp.title,
	cp.content, 
	cp.user_id,
	u.username,
	cp.date
	FROM communities_posts cp
	JOIN users u ON cp.user_id = u.id
	WHERE cp.community_id = $1
	ORDER BY cp.date DESC`

	var posts []*model.CommunityPostResponse
	err := db.db.Select(&posts, query, communityId)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*model.CommunityPostResponse{}, nil
		}
		return nil, fmt.Errorf("failed to get community posts: %w", err)
	}
	return posts, nil
}

func (db *PostgresCommunitiesRepository) CreateCommunityPost(communityId uuid.UUID, userId uuid.UUID, content string, title string) error {
	query := `INSERT INTO communities_posts (community_id, user_id, content, title) VALUES ($1, $2, $3, $4)`
	_, err := db.db.Exec(query, communityId, userId, content, title)
	if err != nil {
		return fmt.Errorf("failed to create community post: %w", err)
	}
	return nil
}

func (db *PostgresCommunitiesRepository) LeaveCommunity(communityId uuid.UUID, userId uuid.UUID) error {
	query := `DELETE FROM communities_users WHERE community_id = $1 AND user_id = $2`
	_, err := db.db.Exec(query, communityId, userId)
	if err != nil {
		return fmt.Errorf("failed to leave community: %w", err)
	}
	return nil
}

func (db *PostgresCommunitiesRepository) CheckIFCommunityExists(communityId uuid.UUID) bool {
	query := `SELECT EXISTS(SELECT 1 FROM communities WHERE id=$1)`

	var exists bool
	err := db.db.QueryRow(query, communityId).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func (db *PostgresCommunitiesRepository) DeleteCommunity(communityId uuid.UUID) error {
	query := `DELETE FROM communities_posts WHERE community_id = $1`
	_, err := db.db.Exec(query, communityId)
	if err != nil {
		return fmt.Errorf("failed to delete community: %w", err)
	}

	query = `DELETE FROM communities_pictures WHERE community_id = $1`
	_, err = db.db.Exec(query, communityId)
	if err != nil {
		return fmt.Errorf("failed to delete community: %w", err)
	}

	// Deletes from users
	query = `DELETE FROM communities_users WHERE community_id = $1`
	_, err = db.db.Exec(query, communityId)
	if err != nil {
		return fmt.Errorf("failed to delete community: %w", err)
	}

	query = `DELETE FROM communities WHERE id = $1`
	_, err = db.db.Exec(query, communityId)
	if err != nil {
		return fmt.Errorf("failed to delete community: %w", err)
	}

	return nil
}

func (db *PostgresCommunitiesRepository) CheckIfUserIsCreator(communityId uuid.UUID, userId uuid.UUID) bool {
	query := `SELECT EXISTS(SELECT 1 FROM communities WHERE id=$1 AND owner_id=$2)`

	var exists bool
	err := db.db.QueryRow(query, communityId, userId).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}
