package usecase_admin

import (
	"app/domain"
	mongo_model "app/domain/model/mongo"
	request_model "app/domain/model/request"
	jwt_helper "app/helpers/jsonwebtoken"
	"context"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"golang.org/x/crypto/bcrypt"
)

func (u *appUsecase) Login(ctx context.Context, payload request_model.LoginRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Email == "" {
		errValidation["email"] = "email field is required"
	}

	if payload.Password == "" {
		errValidation["password"] = "password field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check the db
	admin, err := u.mongodbRepo.FetchOneAdmin(ctx, mongo_model.AdminFilter{
		Email: &payload.Email,
	})
	if err != nil {
		return response.Error(500, err.Error())
	}

	if admin == nil {
		return response.Error(400, "admin not found")
	}

	// check password
	if err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(payload.Password)); err != nil {
		return response.Error(400, "Wrong password")
	}

	// generate token
	tokenString, err := jwt_helper.GenerateJWTToken(
		jwt_helper.GetJwtCredential().Admin,
		domain.JWTClaimAdmin{
			AdminID: admin.ID.Hex(),
		},
	)
	if err != nil {
		return response.Error(400, err.Error())
	}

	return response.Success(map[string]interface{}{
		"admin": admin,
		"token": tokenString,
	})
}
