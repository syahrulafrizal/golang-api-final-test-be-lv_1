package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt" json:"-"`
}

var UserAllowedSort = []string{"name", "username", "createdAt", "updatedAt"}

type UserFilter struct {
	DefaultFilter
	Username *string
}

func (f *UserFilter) Query(defaultQuery map[string]any) map[string]any {
	// default query
	f.DefaultFilter.DefaultQuery(defaultQuery)

	if f.Username != nil {
		defaultQuery["username"] = f.Username
	}

	return defaultQuery
}
