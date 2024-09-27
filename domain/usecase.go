package domain

import (
	request_model "app/domain/model/request"
	"context"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

type MemberAppUsecase interface {
	Login(ctx context.Context, payload request_model.LoginRequest) response.Base
	Register(ctx context.Context, payload request_model.RegisterRequest) response.Base
	GetMe(ctx context.Context, claim JWTClaimUser) response.Base
}
