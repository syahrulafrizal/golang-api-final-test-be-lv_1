package mongo_model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Faq struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt *time.Time         `bson:"deletedAt" json:"-"`
}

var FaqAllowedSort = []string{"title", "createdAt", "updatedAt"}

type FaqFilter struct {
	DefaultFilter
}

func (f *FaqFilter) Query(defaultQuery map[string]any) map[string]any {
	// default query
	f.DefaultFilter.DefaultQuery(defaultQuery)

	return defaultQuery
}
