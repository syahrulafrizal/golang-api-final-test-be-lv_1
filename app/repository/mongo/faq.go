package mongorepo

import (
	mongo_model "app/domain/model/mongo"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// for default query
func defaultFaqQuery() map[string]any {
	return map[string]any{
		"deletedAt": map[string]any{
			"$eq": nil,
		},
	}
}

func (r *mongoDBRepo) FetchFaq(ctx context.Context, options mongo_model.FaqFilter) (cur *mongo.Cursor, err error) {
	// generate query
	query := options.Query(defaultFaqQuery())

	// query options
	findOptions := options.FindOptions()

	cur, err = r.conn.Collection(r.faqCollection).Find(ctx, query, findOptions)
	if err != nil {
		logrus.Error("FetchFaq Find:", err)
		return
	}

	return
}

func (r *mongoDBRepo) FetchOneFaq(ctx context.Context, options mongo_model.FaqFilter) (row *mongo_model.Faq, err error) {
	// generate query
	query := options.Query(defaultFaqQuery())

	err = r.conn.Collection(r.faqCollection).FindOne(ctx, query).Decode(&row)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			return
		}

		logrus.Error("FetchOneFaq FindOne:", err)
		return
	}

	return
}

func (r *mongoDBRepo) CountFaq(ctx context.Context, options mongo_model.FaqFilter) (total int64) {
	// generate query
	query := options.Query(defaultFaqQuery())

	total, err := r.conn.Collection(r.faqCollection).CountDocuments(ctx, query)
	if err != nil {
		logrus.Error("CountFaq", err)
		return 0
	}

	return
}

func (r *mongoDBRepo) CreateFaq(ctx context.Context, row *mongo_model.Faq) (err error) {
	_, err = r.conn.Collection(r.faqCollection).InsertOne(ctx, row)
	if err != nil {
		logrus.Error("CreateFaq InsertOne:", err)
		return
	}
	return
}
