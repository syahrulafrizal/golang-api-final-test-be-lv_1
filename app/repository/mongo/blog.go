package mongorepo

import (
	mongo_model "app/domain/model/mongo"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// for default query
func defaultBlogQuery() map[string]any {
	return map[string]any{
		"deletedAt": map[string]any{
			"$eq": nil,
		},
	}
}

func (r *mongoDBRepo) FetchBlog(ctx context.Context, options mongo_model.BlogFilter) (cur *mongo.Cursor, err error) {
	// generate query
	query := options.Query(defaultBlogQuery())

	// query options
	findOptions := options.FindOptions()

	cur, err = r.conn.Collection(r.blogCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchBlog Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) FetchOneBlog(ctx context.Context, options mongo_model.BlogFilter) (row *mongo_model.Blog, err error) {
	// generate query
	query := options.Query(defaultBlogQuery())

	err = r.conn.Collection(r.blogCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneBlog FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountBlog(ctx context.Context, options mongo_model.BlogFilter) (total int64) {
	// generate query
	query := options.Query(defaultBlogQuery())

	total, err := r.conn.Collection(r.blogCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountBlog", err)
		return 0
	}

	return
}

func (r *mongoDBRepo) CreateBlog(ctx context.Context, row *mongo_model.Blog) (err error) {
	_, err = r.conn.Collection(r.blogCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateBlog InsertOne:", err)
		return
	}
	return
}
