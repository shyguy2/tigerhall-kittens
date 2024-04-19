package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"tigerhall-kittens-app/pkg/auth"
	"tigerhall-kittens-app/pkg/models"
	"time"
)

// mockTigerService is a mock implementation of the TigerService interface.
type mockTigerService struct {
	signupService                func(user *models.User) error
	loginService                 func(credentials models.LoginCredentials) (*models.User, error)
	createTigerService           func(tiger models.Tiger) error
	getAllTigersService          func(page, pageSize int) ([]*models.Tiger, int, error)
	createTigerSighting          func(newSighting *models.TigerSighting) error
	getAllTigerSightings         func(tigerID int) ([]*models.TigerSighting, error)
	createTigerSightingService   func(newSighting *models.TigerSighting) error
	getTigerSightingsByIDService func(tigerID, page, pageSize int) ([]*models.TigerSighting, int, error)
}

func (m *mockTigerService) SignupService(user *models.User) error {
	return m.signupService(user)
}

func (m *mockTigerService) LoginService(credentials models.LoginCredentials) (*models.User, error) {
	return m.loginService(credentials)
}

func (m *mockTigerService) CreateTigerService(tiger models.Tiger) error {
	return m.createTigerService(tiger)
}

func (m *mockTigerService) GetAllTigersService(page, pageSize int) ([]*models.Tiger, int, error) {
	return m.getAllTigersService(page, pageSize)
}

func (m *mockTigerService) CreateTigerSightingService(newSighting *models.TigerSighting) error {
	return m.createTigerSightingService(newSighting)
}

func (m *mockTigerService) GetTigerSightingsByIDService(tigerID, page, pageSize int) ([]*models.TigerSighting, int, error) {
	return m.getTigerSightingsByIDService(tigerID, page, pageSize)
}

func TestSignupHandler_Success(t *testing.T) {
	// Arrange
	user := models.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "testpassword",
	}

	mockService := &mockTigerService{
		signupService: func(user *models.User) error {
			return nil
		},
	}

	handler := NewHandlers(mockService, log.Default(), nil)
	body, _ := json.Marshal(user)
	req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Act
	handler.SignupHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code, "Status code should be 201")
	assert.JSONEq(t, `{"message":"success"}`, rr.Body.String(), "Response body should match")
}

func TestSignupHandler_BadRequest(t *testing.T) {
	// Arrange
	user := models.User{
		// Username and Password fields are intentionally left empty to trigger a bad request error.
		Email: "testuser@example.com",
	}

	mockService := &mockTigerService{}

	handler := NewHandlers(mockService, log.Default(), nil)
	body, _ := json.Marshal(user)
	req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Act
	handler.SignupHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code should be 400")
	assert.Contains(t, rr.Body.String(), "username is required", "Response body should contain error message")
}

func TestSignupHandler_InternalServerError(t *testing.T) {
	// Arrange
	user := models.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "testpassword",
	}

	mockService := &mockTigerService{
		signupService: func(user *models.User) error {
			return errors.New("failed to create user")
		},
	}

	handler := NewHandlers(mockService, log.Default(), nil)
	body, _ := json.Marshal(user)
	req, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Act
	handler.SignupHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code should be 500")
	assert.Contains(t, rr.Body.String(), "failed to create user", "Response body should contain error message")
}

func TestLoginHandler_Success(t *testing.T) {
	// Arrange
	loginCredentials := models.LoginCredentials{
		Email:    "testuser@example.com",
		Password: "testpassword",
	}

	mockService := &mockTigerService{
		loginService: func(credentials models.LoginCredentials) (*models.User, error) {
			// Simulate a successful login and return a user
			return &models.User{
				Username: "testuser",
				Email:    "testuser@example.com",
			}, nil
		},
	}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	body, _ := json.Marshal(loginCredentials)
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Act
	handler.LoginHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.NotEmpty(t, response["token"], "Token should not be empty")
}

func TestLoginHandler_Unauthorized(t *testing.T) {
	// Arrange
	loginCredentials := models.LoginCredentials{
		Email:    "testuser@example.com",
		Password: "wrong_password", // Wrong password to trigger unauthorized error
	}

	mockService := &mockTigerService{
		loginService: func(credentials models.LoginCredentials) (*models.User, error) {
			return nil, errors.New("invalid email or password")
		},
	}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	body, _ := json.Marshal(loginCredentials)
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Act
	handler.LoginHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Status code should be 401")
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.Empty(t, response["token"], "Token should be empty")
}

func TestLoginHandler_BadRequest(t *testing.T) {
	// Arrange
	// Invalid request body (not a valid JSON for loginCredentials)
	invalidBody := []byte(`invalid_json`)

	mockService := &mockTigerService{}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(invalidBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Act
	handler.LoginHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code should be 400")
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.Empty(t, response["token"], "Token should be empty")
}

func TestCreateTigerHandler_Success(t *testing.T) {
	// Arrange
	tiger := models.Tiger{
		Name:        "Mufasa",
		DateOfBirth: time.Date(2015, 1, 15, 0, 0, 0, 0, time.UTC),
		LastSeen:    time.Now(),
		Lat:         12.345,
		Long:        67.890,
	}

	mockService := &mockTigerService{
		createTigerService: func(tiger models.Tiger) error {
			// Simulate a successful tiger creation
			// We can assume that the tiger is added to the database here
			return nil
		},
	}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	body, _ := json.Marshal(tiger)
	req, err := http.NewRequest(http.MethodPost, "/tigers", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Act
	handler.CreateTigerHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code, "Status code should be 201")
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.Equal(t, "success", response["message"], "Message should be 'success'")
}

func TestCreateTigerHandler_BadRequest(t *testing.T) {
	// Arrange
	// Invalid request body (not a valid JSON for tiger)
	invalidBody := []byte(`invalid_json`)

	mockService := &mockTigerService{}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	req, err := http.NewRequest(http.MethodPost, "/tigers", bytes.NewReader(invalidBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Act
	handler.CreateTigerHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code should be 400")
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.NotEmpty(t, response["error"], "Error should not be empty")
}

func TestCreateTigerHandler_InternalServerError(t *testing.T) {
	// Arrange
	tiger := models.Tiger{
		Name:        "Simba",
		DateOfBirth: time.Date(2018, 5, 10, 0, 0, 0, 0, time.UTC),
		LastSeen:    time.Now(),
		Lat:         23.456,
		Long:        87.654,
	}

	mockService := &mockTigerService{
		createTigerService: func(tiger models.Tiger) error {
			// Simulate an error during tiger creation
			return errors.New("failed to create tiger")
		},
	}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	body, _ := json.Marshal(tiger)
	req, err := http.NewRequest(http.MethodPost, "/tigers", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Act
	handler.CreateTigerHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code should be 500")
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.NotEmpty(t, response["error"], "Error should not be empty")
}

func TestGetAllTigersHandler_Success(t *testing.T) {
	// Arrange
	tigers := []*models.Tiger{
		{
			ID:          1,
			Name:        "Mufasa",
			DateOfBirth: time.Now(),
			LastSeen:    time.Now(),
			Lat:         12.345,
			Long:        67.890,
		},
		{
			ID:          2,
			Name:        "Simba",
			DateOfBirth: time.Now(),
			LastSeen:    time.Now(),
			Lat:         23.456,
			Long:        87.654,
		},
	}

	mockService := &mockTigerService{
		getAllTigersService: func(int, int) ([]*models.Tiger, int, error) {
			// Simulate a successful retrieval of tigers
			return tigers, len(tigers), nil
		},
	}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	req, err := http.NewRequest(http.MethodGet, "/tigers?page=1&pageSize=10", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Act
	handler.GetAllTigersHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.Equal(t, float64(2), response["totalCount"], "Expected 2 tigers in response")
}

func TestGetAllTigersHandler_InternalServerError(t *testing.T) {
	// Arrange
	mockService := &mockTigerService{
		getAllTigersService: func(int, int) ([]*models.Tiger, int, error) {
			// Simulate an error during retrieval of tigers
			return nil, 0, errors.New("failed to fetch tigers")
		},
	}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	req, err := http.NewRequest(http.MethodGet, "/tigers?page=1&pageSize=10", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Act
	handler.GetAllTigersHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code should be 500")
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.NotEmpty(t, response["error"], "Error should not be empty")
}

func TestCreateTigerSightingHandler_ValidationError(t *testing.T) {
	// Arrange
	mockService := &mockTigerService{}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	// Create an invalid tiger sighting request (missing required fields)
	tigerSighting := models.TigerSighting{
		TigerID:       1,
		Timestamp:     time.Now(),
		Lat:           12.345,
		Long:          67.890,
		ReporterEmail: "", // Missing reporter email
	}

	// Create a form data payload
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	writer.WriteField("tigerID", strconv.Itoa(tigerSighting.TigerID))
	writer.WriteField("timestamp", tigerSighting.Timestamp.Format(time.RFC3339))
	writer.WriteField("lat", strconv.FormatFloat(tigerSighting.Lat, 'f', -1, 64))
	writer.WriteField("long", strconv.FormatFloat(tigerSighting.Long, 'f', -1, 64))
	writer.WriteField("otherField", "otherValue") // Other unrelated form field
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "/create_tiger_sighting", &requestBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Add the username to the request context for use in the handlers
	ctx := req.Context()
	ctx = context.WithValue(ctx, "username", "rajnish")
	req = req.WithContext(ctx)

	// Add the email to the request context for use in the handlers
	ctx = context.WithValue(ctx, "email", "rajnish.kumar@gmail.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Act
	handler.CreateTigerSightingHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code should be 400")
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.Contains(t, response["error"], "Failed to get image file", "Error should mention missing reporterEmail")
}

func TestCreateTigerSightingHandler_InternalServerError(t *testing.T) {
	// Arrange
	mockService := &mockTigerService{
		createTigerSightingService: func(sighting *models.TigerSighting) error {
			// Simulate an error during creation of tiger sighting
			return errors.New("failed to create tiger sighting")
		},
	}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	// Create a valid tiger sighting request
	tigerSighting := models.TigerSighting{
		TigerID:       1,
		Timestamp:     time.Now(),
		Lat:           12.345,
		Long:          67.890,
		ReporterEmail: "reporter@example.com",
	}

	// Create a form data payload
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	writer.WriteField("tigerID", strconv.Itoa(tigerSighting.TigerID))
	writer.WriteField("timestamp", tigerSighting.Timestamp.Format(time.RFC3339))
	writer.WriteField("lat", strconv.FormatFloat(tigerSighting.Lat, 'f', -1, 64))
	writer.WriteField("long", strconv.FormatFloat(tigerSighting.Long, 'f', -1, 64))
	writer.WriteField("reporterEmail", tigerSighting.ReporterEmail)
	writer.WriteField("otherField", "otherValue") // Other unrelated form field
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "/create_tiger_sighting", &requestBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	// Act
	handler.CreateTigerSightingHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code should be 500")
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.NotEmpty(t, response["error"], "Error should not be empty")
}

func TestGetAllTigerSightingsHandler_Success(t *testing.T) {
	// Arrange
	mockService := &mockTigerService{
		getTigerSightingsByIDService: func(tigerID, page, pageSize int) ([]*models.TigerSighting, int, error) {
			// Simulate a successful retrieval of tiger sightings
			tigerSightings := []*models.TigerSighting{
				{
					ID:        1,
					TigerID:   1,
					Timestamp: time.Now(),
					Lat:       12.345,
					Long:      67.890,
				},
				{
					ID:        2,
					TigerID:   1,
					Timestamp: time.Now(),
					Lat:       12.346,
					Long:      67.891,
				},
			}

			return tigerSightings, len(tigerSightings), nil
		},
	}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	// Prepare a request with "id" query parameter
	req, err := http.NewRequest(http.MethodGet, "tiger/:id/sightings?page=1&pageSize=2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Set the request context to include the query parameters
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Act
	handler.GetTigerSightingsByIDHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")

	var response map[string]interface{}

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.Equal(t, float64(2), response["totalCount"], "Expected 2 tiger sightings in response")
}

func TestGetAllTigerSightingsHandler_InvalidID(t *testing.T) {
	// Arrange
	mockService := &mockTigerService{}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	// Prepare a request with invalid "id" query parameter (not an integer)
	req, err := http.NewRequest(http.MethodGet, "/tiger/:id/sightings?page=1&pageSize=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Set the request context to include the query parameters
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

	// Act
	handler.GetTigerSightingsByIDHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code should be 400")
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.Contains(t, response["error"], "Invalid tiger_id", "Error should mention invalid tiger_id")
}

func TestGetAllTigerSightingsHandler_InternalServerError(t *testing.T) {
	// Arrange
	mockService := &mockTigerService{
		getTigerSightingsByIDService: func(tigerID, page, pageSize int) ([]*models.TigerSighting, int, error) {
			// Simulate an error during retrieval of tiger sightings
			return nil, 0, errors.New("failed to retrieve tiger sightings")
		},
	}

	auth := auth.NewAuth("test_secret_key")
	handler := NewHandlers(mockService, log.Default(), auth)

	// Prepare a request with "id" query parameter
	req, err := http.NewRequest(http.MethodGet, "/tiger/:id/sightings?page=1&pageSize=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Set the request context to include the query parameters
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Act
	handler.GetTigerSightingsByIDHandler(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code should be 500")
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Error while unmarshaling response")
	assert.NotEmpty(t, response["error"], "Error should not be empty")
}
