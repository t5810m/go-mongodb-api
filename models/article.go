package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Article struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string        `bson:"title" json:"title" validate:"required,min=2,max=255"`
	Content     string        `bson:"content" json:"content" validate:"required"`
	Slug        string        `bson:"slug" json:"slug" validate:"required"`
	Active      bool          `bson:"active" json:"active"`
	CreatedTime time.Time     `bson:"created_time" json:"created_time"`
	UpdatedTime time.Time     `bson:"updated_time" json:"updated_time"`
	CreatedBy   string        `bson:"created_by" json:"created_by"`
	UpdatedBy   string        `bson:"updated_by" json:"updated_by"`
}
