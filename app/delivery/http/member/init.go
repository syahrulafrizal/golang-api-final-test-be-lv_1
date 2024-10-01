package http_member

import (
	"app/app/delivery/http/middleware"
	"app/domain"

	"github.com/gin-gonic/gin"
)

type routeHandler struct {
	Usecase    domain.MemberAppUsecase
	Route      *gin.RouterGroup
	Middleware middleware.Middleware
}

func NewRouteHandler(route *gin.RouterGroup, middleware middleware.Middleware, u domain.MemberAppUsecase) {
	handler := &routeHandler{
		Usecase:    u,
		Route:      route,
		Middleware: middleware,
	}

	handler.handleAuthRoute("/auth")
	handler.handleSampleRoute("/sample")
}
