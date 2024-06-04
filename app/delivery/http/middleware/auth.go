package middleware

import (
	"app/domain"
	"errors"
	"net/http"
	"strings"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (m *AppMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		hAuth := c.GetHeader("Authorization")
		if hAuth == "" {
			response := response.Error(http.StatusBadRequest, "Unauthorized: Header authorization is required")
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		splitToken := strings.Split(hAuth, "Bearer ")
		if len(splitToken) != 2 {
			response := response.Error(http.StatusBadRequest, "Unauthorized: Token is invalid")
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// get token without 'Bearer '
		tokenString := splitToken[1]

		// validating token
		token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaimUser{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secret), nil
		})

		// check validity token
		if !token.Valid {
			msg := err.Error()
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				msg = "Unauthorized: Token signature invalid"
			}
			response := response.Error(http.StatusBadRequest, msg)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		claims, tokenOK := token.Claims.(*domain.JWTClaimUser)
		if !tokenOK {
			response := response.Error(http.StatusBadRequest, "Unauthorized: Token data not valid")
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		if !token.Valid {
			response := response.Error(http.StatusBadRequest, err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		c.Set("token_data", *claims)
		c.Next()
	}
}
