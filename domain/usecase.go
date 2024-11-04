package domain

import (
	request_model "app/domain/model/request"
	"context"
	"net/url"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

type AdminAppUsecase interface {
	Login(ctx context.Context, payload request_model.LoginRequest) response.Base
}

type PublicAppUsecase interface {
	FaqList(ctx context.Context, query url.Values) response.Base
}
