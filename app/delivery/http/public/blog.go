package http_public

import (
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleBlogRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.BlogList)
}

func (r *routeHandler) BlogList(c *gin.Context) {
	ctx := c.Request.Context()

	response := r.Usecase.BlogList(ctx, c.Request.URL.Query())
	c.JSON(response.Status, response)
}
