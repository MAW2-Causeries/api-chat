package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// VerifyBearerToken verifies the provided bearer token and returns the user ID if valid
func VerifyBearerToken(Authorization string) (string, error) {
	if Authorization == "" {
		return "", errors.New("missing token")
	}
	tokenString := Authorization[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(GetEnv("JWT_SECRET", "1234")), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	userID, _ := claims["user_id"].(string)
	return userID, nil
}