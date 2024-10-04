package http_member

import (
	"app/domain"

	"github.com/gin-gonic/gin"
)

func (h *routeHandler) handleSampleRoute(prefixPath string) {
	// (optional). add prefix api version
	api := h.Route.Group(prefixPath)

	api.GET("/user/list", h.Middleware.Auth(), h.UserList)
	api.GET("/user/detail/:id", h.Middleware.Auth(), h.UserDetail)
	api.GET("/user/export", h.UserExport)
}

func (r *routeHandler) UserList(c *gin.Context) {
	ctx := c.Request.Context()

	response := r.Usecase.SampleUserList(ctx, c.MustGet("token_data").(domain.JWTClaimUser), c.Request.URL.Query())
	c.JSON(response.Status, response)
}

func (r *routeHandler) UserDetail(c *gin.Context) {
	ctx := c.Request.Context()

	response := r.Usecase.SampleUserDetail(ctx, c.MustGet("token_data").(domain.JWTClaimUser), c.Param("id"))
	c.JSON(response.Status, response)
}

func (r *routeHandler) UserExport(c *gin.Context) {
	ctx := c.Request.Context()

	response := r.Usecase.SampleUserExport(ctx, c.MustGet("token_data").(domain.JWTClaimUser), c.Request.URL.Query())
	c.JSON(response.Status, response)
}
