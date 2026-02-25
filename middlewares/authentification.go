package middlewares

import (
	"MessagesService/utils"
	"context"
	"net/http"
)

// https://go.dev/blog/context#TOC_3.2.
type contextKey string

// UserIDKey is the key used to store the user ID in the request context.
const UserIDKey contextKey = "userID"

// Authentification is a middleware that checks for the presence of a valid Bearer token in the Authorization header.
func Authentification(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Authorization := r.Header.Get("Authorization")
		userID, err := utils.VerifyBearerToken(Authorization)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
