package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"tigerhall-kittens-app/pkg/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
	secretKey string
}

func NewAuth(secretKey string) *Auth {
	return &Auth{
		secretKey: secretKey,
	}
}

func (a *Auth) GenerateToken(username, email string) (string, error) {
	// Create a new token object, specifying signing method and claims
	token := jwt.New(jwt.SigningMethodHS256)

	// Create claims for the token (e.g., username, expiration time)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}

func (a *Auth) VerifyToken(tokenString string) (string, string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for validation
		return []byte(a.secretKey), nil
	})

	if err != nil {
		return "", "", fmt.Errorf("failed to parse token: %v", err)
	}

	// Validate the token and extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", fmt.Errorf("invalid token")
	}

	// Extract and return the username from claims
	username, ok := claims["username"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid token claims")
	}

	// Extract and return the username from claims
	email, ok := claims["email"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid token claims")
	}

	return username, email, nil
}

func GetEmailFromContext(ctx context.Context) (string, bool) {
	// Retrieve the email value from the context
	email, ok := ctx.Value("email").(string)

	return email, ok
}

func ValidateUserData(user models.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	} else if user.Email == "" {
		return errors.New("email is required")
	} else if user.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func HashPassword(password string) (string, error) {
	// Generate the hash of the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Convert the hashed password to a string and return it
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, plainPassword string) error {
	// Compare the hashed password with the plain-text password using bcrypt.CompareHashAndPassword
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		// If the comparison fails, return an error indicating an invalid password
		return err
	}

	return nil
}
