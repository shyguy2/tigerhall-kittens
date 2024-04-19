package middleware

import (
	"context"
	"net/http"
	"strings"
	"tigerhall-kittens-app/pkg/auth"
)

func AuthMiddleware(auth *auth.Auth, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Expect the token to be in the format: "Bearer <token>"
		tokenParts := strings.Split(authorizationHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]

		// Verify the token
		username, email, err := auth.VerifyToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add the username to the request context for use in the handlers
		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", username)
		r = r.WithContext(ctx)

		// Add the email to the request context for use in the handlers
		ctx = context.WithValue(ctx, "email", email)
		r = r.WithContext(ctx)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
