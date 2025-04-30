package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateToken(userID uuid.UUID) (string, error) {

	var secretKey = []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		return "", errors.New("JWT_SECRET not set in .env")
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  userID.String(),
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {

	var secretKey = []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		return "", errors.New("JWT_SECRET not set in .env")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil {
		return "", errors.New("invalid claims")
	}

	if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
		return "", errors.New("token expired")
	}

	return claims["id"].(string), nil
}
