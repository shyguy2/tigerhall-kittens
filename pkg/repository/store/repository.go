package store

import (
	"database/sql"
	"fmt"
	"tigerhall-kittens-app/pkg/models"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *postgresRepository {
	return &postgresRepository{db: db}
}

func (p *postgresRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
	`
	_, err := p.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) CreateTiger(tiger *models.Tiger) error {
	query := `
		INSERT INTO tigers (name, date_of_birth, last_seen, lat, long)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := p.db.Exec(query, tiger.Name, tiger.DateOfBirth, tiger.LastSeen, tiger.Lat, tiger.Long)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) GetAllTigersWithPagination(page, pageSize int) ([]*models.Tiger, int, error) {
	query := `
		SELECT id, name, date_of_birth, last_seen, lat, long
		FROM tigers
		ORDER BY last_seen DESC
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * pageSize

	rows, err := p.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tigers := []*models.Tiger{}
	for rows.Next() {
		tiger := &models.Tiger{}
		err := rows.Scan(&tiger.ID, &tiger.Name, &tiger.DateOfBirth, &tiger.LastSeen, &tiger.Lat, &tiger.Long)
		if err != nil {
			return nil, 0, err
		}
		tigers = append(tigers, tiger)
	}

	// Get total count of tigers (without pagination)
	totalCount, err := p.GetTotalTigerCount()
	if err != nil {
		return nil, 0, err
	}

	return tigers, totalCount, nil
}

func (p *postgresRepository) GetTotalTigerCount() (int, error) {
	query := `
		SELECT COUNT(*) FROM tigers
	`

	var totalCount int
	err := p.db.QueryRow(query).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (p *postgresRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
        SELECT id, username, email, password
        FROM users
        WHERE email = $1
    `

	user := &models.User{}
	err := p.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (p *postgresRepository) CreateTigerSighting(tigerSighting *models.TigerSighting) error {
	query := `
       INSERT INTO tiger_sightings (tiger_id, timestamp, lat, long, image, reporter_Email)
       VALUES ($1, $2, $3, $4, $5,$6)
       RETURNING id
   `
	err := p.db.QueryRow(query, tigerSighting.TigerID, tigerSighting.Timestamp, tigerSighting.Lat, tigerSighting.Long, tigerSighting.Image, tigerSighting.ReporterEmail).Scan(&tigerSighting.ID)
	if err != nil {
		return fmt.Errorf("failed to create tiger sighting: %v", err)
	}

	return nil
}

func (p *postgresRepository) GetTigerSightingsByID(tigerID int) ([]*models.TigerSighting, error) {
	query := "SELECT id, tiger_id, timestamp, lat, long, image, reporter_Email FROM tiger_sightings WHERE tiger_id = $1 ORDER BY timestamp DESC"

	rows, err := p.db.Query(query, tigerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tiger sightings: %v", err)
	}
	defer rows.Close()

	sightings := []*models.TigerSighting{}
	for rows.Next() {
		var sighting models.TigerSighting
		err := rows.Scan(&sighting.ID, &sighting.TigerID, &sighting.Timestamp, &sighting.Lat, &sighting.Long, &sighting.Image, &sighting.ReporterEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tiger sighting: %v", err)
		}
		sightings = append(sightings, &sighting)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error processing tiger sightings rows: %v", err)
	}

	return sightings, nil
}

func (p *postgresRepository) GetTigerSightingsByIDWithPagination(tigerID, page, pageSize int) ([]*models.TigerSighting, int, error) {
	query := `
		SELECT id, tiger_id, timestamp, lat, long, image, reporter_Email
		FROM tiger_sightings
		WHERE tiger_id = $1
		ORDER BY timestamp DESC
		OFFSET $2
		LIMIT $3
	`

	offset := (page - 1) * pageSize

	rows, err := p.db.Query(query, tigerID, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tiger sightings: %v", err)
	}
	defer rows.Close()

	sightings := []*models.TigerSighting{}
	for rows.Next() {
		var sighting models.TigerSighting
		err := rows.Scan(&sighting.ID, &sighting.TigerID, &sighting.Timestamp, &sighting.Lat, &sighting.Long, &sighting.Image, &sighting.ReporterEmail)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan tiger sighting: %v", err)
		}
		sightings = append(sightings, &sighting)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error processing tiger sightings rows: %v", err)
	}

	// Get total count of tiger sightings (without pagination)
	totalCount, err := p.GetTigerSightingsCountByID(tigerID)
	if err != nil {
		return nil, 0, err
	}

	return sightings, totalCount, nil
}

func (p *postgresRepository) GetTigerSightingsCountByID(tigerID int) (int, error) {
	query := `
		SELECT COUNT(*) FROM tiger_sightings WHERE tiger_id = $1
	`

	var totalCount int
	err := p.db.QueryRow(query, tigerID).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (p *postgresRepository) GetPreviousTigerSighting(tigerID int) (*models.TigerSighting, error) {
	// Query the database to get the previous tiger sighting based on tigerID
	query := `
		SELECT id, tiger_id, timestamp, lat, long, image, reporter_Email
		FROM tiger_sightings
		WHERE tiger_id = $1
		ORDER BY timestamp DESC
		LIMIT 1;
	`

	var previousSighting models.TigerSighting
	err := p.db.QueryRow(query, tigerID).Scan(
		&previousSighting.ID,
		&previousSighting.TigerID,
		&previousSighting.Timestamp,
		&previousSighting.Lat,
		&previousSighting.Long,
		&previousSighting.Image,
		&previousSighting.ReporterEmail,
	)

	if err == sql.ErrNoRows {
		// No previous sighting found for the given tigerID
		return nil, nil
	} else if err != nil {
		// Some other error occurred during the query
		return nil, err
	}

	return &previousSighting, nil
}
