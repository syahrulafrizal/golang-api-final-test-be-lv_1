package usecase_admin

import (
	"app/domain"
	"time"
)

type appUsecase struct {
	mongodbRepo    domain.MongoDBRepo
	contextTimeout time.Duration
}

type RepoInjection struct {
	MongoDBRepo domain.MongoDBRepo
}

func NewAppUsecase(r RepoInjection, timeout time.Duration) domain.AdminAppUsecase {
	return &appUsecase{
		mongodbRepo:    r.MongoDBRepo,
		contextTimeout: timeout,
	}
}
