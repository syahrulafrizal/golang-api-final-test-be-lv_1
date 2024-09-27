package mongorepo

import (
	"app/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDBRepo struct {
	conn           *mongo.Database
	userCollection string
}

func NewMongodbRepo(Conn *mongo.Database) domain.MongoDBRepo {
	return &mongoDBRepo{
		conn:           Conn,
		userCollection: "users",
	}
}
