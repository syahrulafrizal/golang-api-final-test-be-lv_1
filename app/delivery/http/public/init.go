package http_public

import (
	"app/app/delivery/http/middleware"
	"app/domain"

	"github.com/gin-gonic/gin"
)

type routeHandler struct {
	Usecase    domain.PublicAppUsecase
	Route      *gin.RouterGroup
	Middleware middleware.Middleware
}

func NewRouteHandler(route *gin.RouterGroup, middleware middleware.Middleware, u domain.PublicAppUsecase) {
	handler := &routeHandler{
		Usecase:    u,
		Route:      route,
		Middleware: middleware,
	}

	handler.handleFaqRoute("/faq")
	handler.handleBlogRoute("/blog")
}
