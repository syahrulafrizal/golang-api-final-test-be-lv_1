package http_admin

import (
	"app/domain"
	request_model "app/domain/model/request"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleBlogRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.Middleware.Auth(), h.BlogList)
	api.POST("/create", h.Middleware.Auth(), h.BlogCreate)
}

func (r *routeHandler) BlogList(c *gin.Context) {
	ctx := c.Request.Context()

	response := r.Usecase.BlogList(ctx, c.MustGet("token_data").(domain.JWTClaimAdmin), c.Request.URL.Query())
	c.JSON(response.Status, response)
}

func (r *routeHandler) BlogCreate(c *gin.Context) {
	ctx := c.Request.Context()

	payload := request_model.BlogRequest{}
	c.Bind(&payload)

	response := r.Usecase.BlogCreate(ctx, c.MustGet("token_data").(domain.JWTClaimAdmin), payload, c.Request)
	c.JSON(response.Status, response)
}
