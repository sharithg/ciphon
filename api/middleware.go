package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *Application) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, errors.New("authorization header not found"))
			return
		}

		// Split the Authorization header to extract the token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedErrorResponse(w, r, errors.New("invalid authorization header format"))
			return
		}

		tokenStr := parts[1]
		claims, err := app.Auth.VerifyToken(tokenStr)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, errors.New("invalid token"))
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims.UserId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
