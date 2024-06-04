package http

import (
	"app/app/delivery/http/middleware"
	"app/app/usecase"

	"github.com/gin-gonic/gin"
)

type routeHandler struct {
	Usecase usecase.AppUsecase
	Route   *gin.Engine
}

func NewRouteHandler(ginEngine *gin.Engine, u usecase.AppUsecase) {
	middle := middleware.NewMiddleware()
	handler := &routeHandler{
		Usecase: u,
		Route:   ginEngine,
	}

	handler.handleAuthRoute("/auth", middle)
}
