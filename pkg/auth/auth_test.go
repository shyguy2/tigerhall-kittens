package auth

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"tigerhall-kittens-app/pkg/models"
)

func TestGenerateToken(t *testing.T) {
	secretKey := "test-secret-key"
	auth := NewAuth(secretKey)

	username := "testuser"
	email := "test@example.com"
	expectedExp := time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := auth.GenerateToken(username, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token and verify the claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, username, claims["username"])
	assert.Equal(t, email, claims["email"])
	assert.Equal(t, float64(expectedExp), claims["exp"])
}

func TestVerifyToken(t *testing.T) {
	secretKey := "test-secret-key"
	auth := NewAuth(secretKey)

	username := "testuser"
	email := "test@example.com"

	tokenString, err := auth.GenerateToken(username, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	verifiedUsername, verifiedEmail, err := auth.VerifyToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, username, verifiedUsername)
	assert.Equal(t, email, verifiedEmail)
}

func TestGetEmailFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "email", "test@example.com")

	email, ok := GetEmailFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "test@example.com", email)
}

func TestValidateUserData(t *testing.T) {
	// Valid user data
	validUser := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
	}

	err := ValidateUserData(validUser)
	assert.NoError(t, err)

	// Invalid user data with missing fields
	invalidUser := models.User{
		Username: "testuser",
		Password: "testpassword",
	}

	err = ValidateUserData(invalidUser)
	assert.Error(t, err)
	assert.Equal(t, "email is required", err.Error())
}

func TestHashPasswordAndVerifyPassword(t *testing.T) {
	plainPassword := "testpassword"

	// Test HashPassword
	hashedPassword, err := HashPassword(plainPassword)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Test VerifyPassword
	err = VerifyPassword(hashedPassword, plainPassword)
	assert.NoError(t, err)

	// Test VerifyPassword with wrong password
	err = VerifyPassword(hashedPassword, "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)
}
