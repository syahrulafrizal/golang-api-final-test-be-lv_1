package domain

import (
	request_model "app/domain/model/request"
	"context"
	"net/http"
	"net/url"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

type AdminAppUsecase interface {
	Login(ctx context.Context, payload request_model.LoginRequest) response.Base

	BlogList(ctx context.Context, claim JWTClaimAdmin, query url.Values) response.Base
	BlogCreate(ctx context.Context, claim JWTClaimAdmin, payload request_model.BlogRequest, request *http.Request) response.Base
}

type PublicAppUsecase interface {
	FaqList(ctx context.Context, query url.Values) response.Base

	BlogList(ctx context.Context, query url.Values) response.Base
}
