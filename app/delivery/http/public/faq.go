package http_public

import (
	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleAuthRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/list", h.FaqList)
}

func (r *routeHandler) FaqList(c *gin.Context) {
	ctx := c.Request.Context()

	response := r.Usecase.FaqList(ctx, c.Request.URL.Query())
	c.JSON(response.Status, response)
}
