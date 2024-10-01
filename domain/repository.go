package domain

import (
	mongo_model "app/domain/model/mongo"
	storage_model "app/domain/model/storage"
	"context"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBRepo interface {
	FetchUser(ctx context.Context, options mongo_model.UserFilter) (*mongo.Cursor, error)
	FetchOneUser(ctx context.Context, options mongo_model.UserFilter) (*mongo_model.User, error)
	CountUser(ctx context.Context, options mongo_model.UserFilter) int64
	CreateUser(ctx context.Context, model *mongo_model.User) (err error)
}

type CacheRepo interface {
	Enabled() bool
	GetTTL() time.Duration
	Get(ctx context.Context, key string) (value []byte, err error)
	Set(ctx context.Context, key string, value []byte, expiration *time.Duration) (err error)
}

type StorageRepo interface {
	GetPresignedLink(objectKey string, expires *time.Duration) string
	GetPublicLink(objectKey string) string
	UploadFilePublic(objectKey string, body io.Reader, contentType string) (uploadData *storage_model.UploadResponse, err error)
	UploadFilePrivate(objectKey string, body io.Reader, contentType string, expires *time.Duration) (uploadData *storage_model.UploadResponse, err error)
}
