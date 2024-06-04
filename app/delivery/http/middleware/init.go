package middleware

import (
	"app/helpers"
)

type AppMiddleware struct {
	secret string
}

func NewMiddleware() *AppMiddleware {
	return &AppMiddleware{
		secret: helpers.GetJWTSecretKey(),
	}
}
