package mongorepo

import (
	mongo_model "app/domain/model/mongo"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// for default query
func defaultUserQuery() map[string]any {
	return map[string]any{
		"deletedAt": map[string]any{
			"$eq": nil,
		},
	}
}

func (r *mongoDBRepo) FetchUser(ctx context.Context, options mongo_model.UserFilter) (cur *mongo.Cursor, err error) {
	// generate query
	query := options.Query(defaultUserQuery())

	// query options
	findOptions := options.FindOptions()

	cur, err = r.conn.Collection(r.userCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchUser Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) FetchOneUser(ctx context.Context, options mongo_model.UserFilter) (row *mongo_model.User, err error) {
	// generate query
	query := options.Query(defaultUserQuery())

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

func (r *mongoDBRepo) CountUser(ctx context.Context, options mongo_model.UserFilter) (total int64) {
	// generate query
	query := options.Query(defaultUserQuery())

	total, err := r.conn.Collection(r.userCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountUser", err)
		return 0
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
