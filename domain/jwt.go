package domain

import "github.com/golang-jwt/jwt/v5"

type JWTClaimAdmin struct {
	AdminID string `json:"adminID"`
	jwt.RegisteredClaims
}
