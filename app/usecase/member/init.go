package usecase_member

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

func NewAppUsecase(r RepoInjection, timeout time.Duration) domain.MemberAppUsecase {
	return &appUsecase{
		mongodbRepo:    r.MongoDBRepo,
		contextTimeout: timeout,
	}
}
