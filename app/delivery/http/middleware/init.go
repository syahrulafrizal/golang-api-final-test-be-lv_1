package middleware

import (
	"app/helpers"

	"github.com/gin-gonic/gin"
)

type appMiddleware struct {
	secret string
}

func NewMiddleware() Middleware {
	return &appMiddleware{
		secret: helpers.GetJWTSecretKey(),
	}
}

type Middleware interface {
	Auth() gin.HandlerFunc
}
