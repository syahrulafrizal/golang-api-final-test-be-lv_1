package usecase_member

import (
	"app/domain"
	mongo_model "app/domain/model/mongo"
	request_model "app/domain/model/request"
	jwt_helper "app/helpers/jsonwebtoken"
	"context"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (u *appUsecase) Login(ctx context.Context, payload request_model.LoginRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Username == "" {
		errValidation["username"] = "username field is required"
	}

	if payload.Password == "" {
		errValidation["password"] = "password field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check the db
	user, err := u.mongodbRepo.FetchOneUser(ctx, mongo_model.UserFilter{
		Username: &payload.Username,
	})
	if err != nil {
		return response.Error(500, err.Error())
	}

	if user == nil {
		return response.Error(400, "user not found")
	}

	// check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return response.Error(400, "Wrong password")
	}

	// generate token
	tokenString, err := jwt_helper.GenerateJWTToken(
		jwt_helper.GetJwtCredential().Member,
		domain.JWTClaimUser{
			UserID: user.ID.Hex(),
		},
	)
	if err != nil {
		return response.Error(400, err.Error())
	}

	return response.Success(map[string]interface{}{
		"user":  user,
		"token": tokenString,
	})
}

func (u *appUsecase) Register(ctx context.Context, payload request_model.RegisterRequest) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)
	// validating request
	if payload.Name == "" {
		errValidation["name"] = "name field is required"
	}

	if payload.Username == "" {
		errValidation["username"] = "username field is required"
	}

	if payload.Password == "" {
		errValidation["password"] = "password field is required"
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	// check the db
	user, err := u.mongodbRepo.FetchOneUser(ctx, mongo_model.UserFilter{
		Username: &payload.Username,
	})
	if err != nil {
		return response.Error(500, err.Error())
	}

	if user != nil {
		return response.Error(400, "username already taken")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	newUser := mongo_model.User{
		ID:        primitive.NewObjectID(),
		Name:      payload.Name,
		Username:  payload.Username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.mongodbRepo.CreateUser(ctx, &newUser)
	if err != nil {
		return response.Error(500, err.Error())
	}

	return response.Success(newUser)
}

func (u *appUsecase) GetMe(ctx context.Context, claim domain.JWTClaimUser) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	userID := claim.UserID

	// check the db
	user, err := u.mongodbRepo.FetchOneUser(ctx, mongo_model.UserFilter{
		DefaultFilter: mongo_model.DefaultFilter{
			IDStr: &userID,
		},
	})
	if err != nil {
		return response.Error(500, err.Error())
	}

	if user == nil {
		return response.Error(400, "user not found")
	}

	return response.Success(user)
}
