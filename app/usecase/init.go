package usecase

import (
	mongorepo "app/app/repository/mongo"
	"context"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

type appUsecase struct {
	mongodbRepo    mongorepo.MongoDBRepo
	contextTimeout time.Duration
}

func NewAppUsecase(r mongorepo.MongoDBRepo, timeout time.Duration) AppUsecase {
	return &appUsecase{
		mongodbRepo:    r,
		contextTimeout: timeout,
	}
}

type AppUsecase interface {
	Login(ctx context.Context, options map[string]interface{}) response.Base
	Register(ctx context.Context, options map[string]interface{}) response.Base
	GetMe(ctx context.Context, options map[string]interface{}) response.Base
}
