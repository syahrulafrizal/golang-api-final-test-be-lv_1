package mongorepo

import (
	mongo_model "app/domain/model/mongo"
	"app/helpers"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

// for default query
var DEFAULT_QUERY_USER = map[string]any{
	"deletedAt": map[string]any{
		"$eq": nil,
	},
}

func generateQueryFilterUser(options map[string]interface{}, withOptions bool) (query bson.M, mongoOptions *moptions.FindOptions) {
	// common filter and find options
	query = helpers.CommonFilter(options, DEFAULT_QUERY_USER)
	if withOptions {
		mongoOptions = helpers.CommonMongoFindOptions(options)
	}

	// your own filter
	if username, ok := options["username"].(string); ok {
		query["username"] = username
	}

	return query, mongoOptions
}

func (r *mongoDBRepo) FetchOneUser(ctx context.Context, options map[string]interface{}) (row *mongo_model.User, err error) {
	query, _ := generateQueryFilterUser(options, false)

	err = r.conn.Collection(r.userCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneUser FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CreateUser(ctx context.Context, row *mongo_model.User) (err error) {
	_, err = r.conn.Collection(r.userCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateUser InsertOne:", err)
		return
	}
	return
}
