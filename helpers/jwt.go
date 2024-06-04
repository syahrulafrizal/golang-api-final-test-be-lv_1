package helpers

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GetJWTSecretKey() string {
	return os.Getenv("JWT_SECRET_KEY")
}

func GenerateJWTToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(GetJWTSecretKey()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
