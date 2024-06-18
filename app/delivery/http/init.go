package http

import (
	"app/app/delivery/http/middleware"
	"app/app/usecase"

	"github.com/gin-gonic/gin"
)

type routeHandler struct {
	Usecase    usecase.AppUsecase
	Route      *gin.Engine
	Middleware middleware.Middleware
}

func NewRouteHandler(ginEngine *gin.Engine, middleware middleware.Middleware, u usecase.AppUsecase) {
	handler := &routeHandler{
		Usecase:    u,
		Route:      ginEngine,
		Middleware: middleware,
	}

	handler.handleAuthRoute("/auth")
}
