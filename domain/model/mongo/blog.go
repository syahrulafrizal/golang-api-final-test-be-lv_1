package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID        primitive.ObjectID  `bson:"_id" json:"id"`
	Title     string              `bson:"title" json:"title"`
	Content   string              `bson:"content" json:"content"`
	Thumbnail string              `bson:"thumbnail" json:"thumbnail"`
	CreatedBy BlogCreatedByNested `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time           `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time          `bson:"deletedAt" json:"-"`
}

type BlogCreatedByNested struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

var BlogAllowedSort = []string{"title", "createdAt", "updatedAt"}

type BlogFilter struct {
	DefaultFilter
}

func (f *BlogFilter) Query(defaultQuery map[string]any) map[string]any {
	// default query
	f.DefaultFilter.DefaultQuery(defaultQuery)

	return defaultQuery
}
