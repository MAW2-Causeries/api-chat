package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthentificationRejectsUnauthorizedRequests(t *testing.T) {
	var called bool
	handler := Authentification(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "Unauthorized\n", rec.Body.String())
}

func TestAuthentificationInjectsUserIDIntoContext(t *testing.T) {
	t.Setenv("JWT_SECRET", "middleware-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "user-456",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	signedToken, err := token.SignedString([]byte("middleware-secret"))
	assert.NoError(t, err)

	var gotUserID any
	handler := Authentification(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID = r.Context().Value(UserIDKey)
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+signedToken)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "user-456", gotUserID)
}
