package usecase

import (
	"app/domain"
	"app/domain/model"
	"app/helpers"
	"context"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (u *appUsecase) Login(ctx context.Context, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// payload
	payload := options["payload"].(domain.LoginRequest)

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
	user, err := u.mongodbRepo.FetchOneUser(ctx, map[string]interface{}{
		"username": payload.Username,
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
	tokenString, err := helpers.GenerateJWTToken(domain.JWTClaimUser{
		UserID: user.ID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		// "userID": "bar",
		// "nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	if err != nil {
		return response.Error(400, err.Error())
	}

	return response.Success(map[string]interface{}{
		"user":  user,
		"token": tokenString,
	})
}

func (u *appUsecase) Register(ctx context.Context, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// payload
	payload := options["payload"].(domain.RegisterRequest)

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
	user, err := u.mongodbRepo.FetchOneUser(ctx, map[string]interface{}{
		"username": payload.Username,
	})
	if err != nil {
		return response.Error(500, err.Error())
	}

	if user != nil {
		return response.Error(400, "username already taken")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	newUser := model.User{
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

func (u *appUsecase) GetMe(ctx context.Context, options map[string]interface{}) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// get claim
	claim := options["claim"].(domain.JWTClaimUser)

	userID := claim.UserID

	// check the db
	user, err := u.mongodbRepo.FetchOneUser(ctx, map[string]interface{}{
		"id": userID,
	})
	if err != nil {
		return response.Error(500, err.Error())
	}

	if user == nil {
		return response.Error(400, "user not found")
	}

	return response.Success(user)
}
