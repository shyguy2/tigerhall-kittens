package server_test

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"tigerhall-kittens-app/pkg/auth"
	"tigerhall-kittens-app/pkg/models"
	"tigerhall-kittens-app/pkg/server"
)

// mockTigerService is a mock implementation of the TigerService interface.
type mockTigerService struct {
	signupService               func(user *models.User) error
	loginService                func(credentials models.LoginCredentials) (*models.User, error)
	createTigerService          func(tiger models.Tiger) error
	getAllTigersService         func() ([]*models.Tiger, error)
	createTigerSightingService  func(newSighting *models.TigerSighting) error
	getAllTigerSightingsService func(tigerID int) ([]*models.TigerSighting, error)
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

func (m *mockTigerService) GetAllTigersService() ([]*models.Tiger, error) {
	return m.getAllTigersService()
}

func (m *mockTigerService) CreateTigerSightingService(newSighting *models.TigerSighting) error {
	return m.createTigerSightingService(newSighting)
}

func (m *mockTigerService) GetAllTigerSightingsService(tigerID int) ([]*models.TigerSighting, error) {
	return m.getAllTigerSightingsService(tigerID)
}

func TestServer_SetupRoutes(t *testing.T) {
	// Arrange
	mockService := &mockTigerService{}
	auth := auth.NewAuth("test_secret_key")
	srv := server.NewServer()
	srv.SetupRoutes(mockService, auth)

	// Act & Assert
	router := mux.NewRouter()

	// Check if the routes are set up correctly
	routes := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		assert.NoError(t, err, "Error getting path template")
		t.Logf("Path: %s", pathTemplate)
		return nil
	})

	assert.NoError(t, routes, "Error walking routes")
}

func TestServer_Start(t *testing.T) {
	// Arrange
	mockService := &mockTigerService{}
	auth := auth.NewAuth("test_secret_key")
	srv := server.NewServer()
	srv.SetupRoutes(mockService, auth)

	// Start the server on a test port
	go func() {
		err := srv.Start("8080")
		assert.NoError(t, err, "Error starting server")
	}()

	// Send a test request to the running server
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	assert.NoError(t, err, "Error creating request")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Error sending request")

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected status code 404")

}
