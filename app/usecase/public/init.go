package usecase_public

import (
	"app/domain"
	"time"
)

type publicUsecase struct {
	mongodbRepo    domain.MongoDBRepo
	contextTimeout time.Duration
}

type RepoInjection struct {
	MongoDBRepo domain.MongoDBRepo
}

func NewAppUsecase(r RepoInjection, timeout time.Duration) domain.PublicAppUsecase {
	return &publicUsecase{
		mongodbRepo:    r.MongoDBRepo,
		contextTimeout: timeout,
	}
}
