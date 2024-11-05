package mongorepo

import (
	"app/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDBRepo struct {
	conn            *mongo.Database
	adminCollection string
	faqCollection   string
}

func NewMongodbRepo(Conn *mongo.Database) domain.MongoDBRepo {
	return &mongoDBRepo{
		conn:            Conn,
		adminCollection: "admins",
		faqCollection:   "faqs",
	}
}
