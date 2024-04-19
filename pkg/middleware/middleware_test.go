package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"tigerhall-kittens-app/pkg/auth"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Create a new instance of Auth with a dummy secret key
	authService := auth.NewAuth("test-secret-key")

	// Generate a valid token
	validToken, err := authService.GenerateToken("testuser", "test@example.com")
	assert.NoError(t, err)

	// Create a new request with the valid token in the Authorization header
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a mock handler that will be called after the AuthMiddleware
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the username and email from the request context
		username := r.Context().Value("username").(string)
		email := r.Context().Value("email").(string)

		// Check if the values in the context match the expected values
		assert.Equal(t, "testuser", username)
		assert.Equal(t, "test@example.com", email)

		// Write a response to indicate that the handler is called
		w.Write([]byte("Handler called"))
	})

	// Call the AuthMiddleware with the mockHandler
	authMiddleware := AuthMiddleware(authService, mockHandler)
	authMiddleware.ServeHTTP(rr, req)

	// Check if the response status code is 200 (OK)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check if the response body is "Handler called"
	assert.Equal(t, "Handler called", rr.Body.String())
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Create a new instance of Auth with a dummy secret key
	authService := auth.NewAuth("test-secret-key")

	// Create a request without the Authorization header (invalid token)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a mock handler that will be called after the AuthMiddleware
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The handler should not be called if the token is invalid
		t.Errorf("Handler should not be called")
	})

	// Call the AuthMiddleware with the mockHandler
	authMiddleware := AuthMiddleware(authService, mockHandler)
	authMiddleware.ServeHTTP(rr, req)

	// Check if the response status code is 401 (Unauthorized)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Check if the response body is "Unauthorized"
	assert.Equal(t, "Unauthorized\n", rr.Body.String())
}
