package store

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"tigerhall-kittens-app/pkg/models"
)

func TestPostgresRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	// Test case data
	username := "testuser"
	email := "testuser@example.com"
	password := "testpassword"
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Expect the INSERT query to be executed
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Username, user.Email, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser(user)
	assert.NoError(t, err)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("failed to meet expectations: %v", err)
	}
}

func TestPostgresRepository_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	// Test case data
	email := "testuser@example.com"
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    email,
		Password: "testpassword",
	}

	// Mock the SELECT query to return the test case data
	mock.ExpectQuery("SELECT id, username, email, password").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(user.ID, user.Username, user.Email, user.Password))

	resultUser, err := repo.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.Equal(t, user, resultUser)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("failed to meet expectations: %v", err)
	}
}

func TestPostgresRepository_CreateTiger(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	// Test case data
	tiger := &models.Tiger{
		Name:        "Tiger 1",
		DateOfBirth: time.Date(2018, 1, 15, 0, 0, 0, 0, time.UTC),
		LastSeen:    time.Now(),
		Lat:         12.3456,
		Long:        78.91011,
	}

	// Mock the INSERT query to return success
	mock.ExpectExec("INSERT INTO tigers").
		WithArgs(tiger.Name, tiger.DateOfBirth, tiger.LastSeen, tiger.Lat, tiger.Long).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateTiger(tiger)
	assert.NoError(t, err)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("failed to meet expectations: %v", err)
	}
}

func TestPostgresRepository_CreateTigerSighting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	// Test case data
	tigerSighting := &models.TigerSighting{
		TigerID:       1,
		Timestamp:     time.Now(),
		Lat:           12.3456,
		Long:          78.91011,
		Image:         []byte("test image data"),
		ReporterEmail: "testuser@example.com",
	}

	// Mock the INSERT query to return the test case data
	mock.ExpectQuery("INSERT INTO tiger_sightings").
		WithArgs(tigerSighting.TigerID, tigerSighting.Timestamp, tigerSighting.Lat, tigerSighting.Long, tigerSighting.Image, tigerSighting.ReporterEmail).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = repo.CreateTigerSighting(tigerSighting)
	assert.NoError(t, err)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("failed to meet expectations: %v", err)
	}
}

func TestPostgresRepository_GetPreviousTigerSighting(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	// Test case data
	tigerID := 1
	tigerSighting := &models.TigerSighting{
		ID:            1,
		TigerID:       tigerID,
		Timestamp:     time.Now(),
		Lat:           12.3456,
		Long:          78.91011,
		Image:         []byte("sample-image"),
		ReporterEmail: "reporter@example.com",
	}

	// Mock the query to return a single row result
	mock.ExpectQuery("SELECT id, tiger_id, timestamp, lat, long, image, reporter_Email FROM tiger_sightings").
		WithArgs(tigerID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "tiger_id", "timestamp", "lat", "long", "image", "reporter_Email"}).
			AddRow(tigerSighting.ID, tigerSighting.TigerID, tigerSighting.Timestamp, tigerSighting.Lat, tigerSighting.Long, tigerSighting.Image, tigerSighting.ReporterEmail))

	// Call the function
	previousSighting, err := repo.GetPreviousTigerSighting(tigerID)

	// Check the result
	assert.NoError(t, err)
	assert.NotNil(t, previousSighting)
	assert.Equal(t, tigerSighting.ID, previousSighting.ID)
	assert.Equal(t, tigerSighting.TigerID, previousSighting.TigerID)
	assert.Equal(t, tigerSighting.Timestamp, previousSighting.Timestamp)
	assert.Equal(t, tigerSighting.Lat, previousSighting.Lat)
	assert.Equal(t, tigerSighting.Long, previousSighting.Long)
	assert.Equal(t, tigerSighting.Image, previousSighting.Image)
	assert.Equal(t, tigerSighting.ReporterEmail, previousSighting.ReporterEmail)

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("failed to meet expectations: %v", err)
	}
}
