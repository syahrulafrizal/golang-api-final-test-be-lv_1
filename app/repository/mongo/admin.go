package mongorepo

import (
	mongo_model "app/domain/model/mongo"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// for default query
func defaultAdminQuery() map[string]any {
	return map[string]any{
		"deletedAt": map[string]any{
			"$eq": nil,
		},
	}
}

func (r *mongoDBRepo) FetchAdmin(ctx context.Context, options mongo_model.AdminFilter) (cur *mongo.Cursor, err error) {
	// generate query
	query := options.Query(defaultAdminQuery())

	// query options
	findOptions := options.FindOptions()

	cur, err = r.conn.Collection(r.adminCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchAdmin Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) FetchOneAdmin(ctx context.Context, options mongo_model.AdminFilter) (row *mongo_model.Admin, err error) {
	// generate query
	query := options.Query(defaultAdminQuery())

	err = r.conn.Collection(r.adminCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneAdmin FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountAdmin(ctx context.Context, options mongo_model.AdminFilter) (total int64) {
	// generate query
	query := options.Query(defaultAdminQuery())

	total, err := r.conn.Collection(r.adminCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountAdmin", err)
		return 0
	}

	return
}

func (r *mongoDBRepo) CreateAdmin(ctx context.Context, row *mongo_model.Admin) (err error) {
	_, err = r.conn.Collection(r.adminCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateAdmin InsertOne:", err)
		return
	}
	return
}
