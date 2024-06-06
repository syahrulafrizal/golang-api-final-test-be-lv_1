package http

import (
	"app/domain"
	"net/http"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleAuthRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.POST("/login", h.Login)
	api.POST("/register", h.Register)

	api.GET("/me", h.Middleware.Auth(), h.GetMe)
}

func (r *routeHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.LoginRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	options := map[string]interface{}{
		"payload": payload,
		"query":   c.Request.URL.Query(),
	}

	response := r.Usecase.Login(ctx, options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.RegisterRequest{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, "invalid json data"))
		return
	}

	options := map[string]interface{}{
		"payload": payload,
		"query":   c.Request.URL.Query(),
	}

	response := r.Usecase.Register(ctx, options)
	c.JSON(response.Status, response)
}

func (r *routeHandler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()

	tokenData := c.MustGet("token_data")

	options := map[string]interface{}{
		"claim": tokenData.(domain.JWTClaimUser),
		"query": c.Request.URL.Query(),
	}

	response := r.Usecase.GetMe(ctx, options)
	c.JSON(response.Status, response)
}
