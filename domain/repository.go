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
	FetchAdmin(ctx context.Context, options mongo_model.AdminFilter) (*mongo.Cursor, error)
	FetchOneAdmin(ctx context.Context, options mongo_model.AdminFilter) (*mongo_model.Admin, error)
	CountAdmin(ctx context.Context, options mongo_model.AdminFilter) int64
	CreateAdmin(ctx context.Context, model *mongo_model.Admin) (err error)

	FetchFaq(ctx context.Context, options mongo_model.FaqFilter) (*mongo.Cursor, error)
	FetchOneFaq(ctx context.Context, options mongo_model.FaqFilter) (*mongo_model.Faq, error)
	CountFaq(ctx context.Context, options mongo_model.FaqFilter) int64
	CreateFaq(ctx context.Context, model *mongo_model.Faq) (err error)

	FetchBlog(ctx context.Context, options mongo_model.BlogFilter) (*mongo.Cursor, error)
	FetchOneBlog(ctx context.Context, options mongo_model.BlogFilter) (*mongo_model.Blog, error)
	CountBlog(ctx context.Context, options mongo_model.BlogFilter) int64
	CreateBlog(ctx context.Context, model *mongo_model.Blog) (err error)
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
