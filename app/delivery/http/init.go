package http

import (
	"app/app/delivery/http/middleware"
	"app/app/usecase"

	"github.com/gin-gonic/gin"
)

type routeHandler struct {
	Usecase    usecase.AppUsecase
	Route      *gin.RouterGroup
	Middleware middleware.Middleware
}

func NewRouteHandler(route *gin.RouterGroup, middleware middleware.Middleware, u usecase.AppUsecase) {
	handler := &routeHandler{
		Usecase:    u,
		Route:      route,
		Middleware: middleware,
	}

	handler.handleAuthRoute("/auth")
}
