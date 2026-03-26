package utils

import (
	"testing"
	"time"

	"github.com/bouk/monkey"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestVerifyBearerToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "user-123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	signedToken, err := token.SignedString([]byte("test-secret"))
	assert.NoError(t, err)

	userID, err := VerifyBearerToken("Bearer " + signedToken)

	assert.NoError(t, err)
	assert.Equal(t, "user-123", userID)
}

func TestVerifyBearerTokenMissingToken(t *testing.T) {
	userID, err := VerifyBearerToken("")

	assert.Error(t, err)
	assert.Equal(t, "missing token", err.Error())
	assert.Empty(t, userID)
}

func TestVerifyBearerTokenInvalidPrefix(t *testing.T) {
	userID, err := VerifyBearerToken("Token abc")

	assert.Error(t, err)
	assert.Equal(t, "invalid token", err.Error())
	assert.Empty(t, userID)
}

func TestVerifyBearerTokenInvalidToken(t *testing.T) {
	userID, err := VerifyBearerToken("Bearer invalid-token")

	assert.Error(t, err)
	assert.Equal(t, "invalid token", err.Error())
	assert.Empty(t, userID)
}

func TestVerifyBearerTokenInvalidClaims(t *testing.T) {
	defer monkey.UnpatchAll()

	monkey.Patch(jwt.Parse, func(_ string, _ jwt.Keyfunc, _ ...jwt.ParserOption) (*jwt.Token, error) {
		return &jwt.Token{
			Valid:  true,
			Claims: jwt.RegisteredClaims{},
		}, nil
	})

	userID, err := VerifyBearerToken("Bearer token")

	assert.Error(t, err)
	assert.Equal(t, "invalid token claims", err.Error())
	assert.Empty(t, userID)
}
