package repository



type PostgresCommunitiesRepository struct {
	db *sql.DB
}

func NewPostgresCommunitiesRepository(db *sql.DB) *PostgresCommunitiesRepository {
	enableUUIDExtension := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	if _, err := c.Exec(enableUUIDExtension); err != nil {
		return nil, fmt.Errorf("failed to enable uuid extension: %w", err)
	}

	schemaCommunities := `
		CREATE TABLE IF NOT EXISTS communities (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			owner_id UUID NOT NULL,
			users UUID NOT NULL,
			posts UUID NOT NULL,

			PRIMARY KEY (id),
			FOREIGN KEY (owner_id) REFERENCES users(id),
			FOREIGN KEY (users) REFERENCES communities_users(id),
			FOREIGN KEY (posts) REFERENCES communities_posts(id)
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

	schemaCommunitiesPosts := `
		CREATE TABLE IF NOT EXISTS communities_posts (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			community_id UUID NOT NULL,
			description TEXT NOT NULL,
			PRIMARY KEY (id),			
			FOREIGN KEY (community_id) REFERENCES communities(id)
		);`
	
	if _, err := db.Exec(schemaCommunitiesPosts); err != nil {
		return nil, fmt.Errorf("failed to create communities_posts table: %w", err)
	}
	
	return &PostgresCommunitiesRepository{db: db}
}

func (db *PostgresCommunitiesRepository) CreateCommunity(community model.NewCommunityRequest) (UUID.uuid, error) {
	// Create a new community in the database
	
	// Prepare the SQL query
	query := `INSERT INTO communities (name, description, owner_id) VALUES ($1, $2, $3) RETURNING id`
	
	// Return the ID of the created community
	return UUID.New(), nil

}