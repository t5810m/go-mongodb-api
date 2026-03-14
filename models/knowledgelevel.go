package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type KnowledgeLevel struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string        `bson:"title" json:"title" validate:"required,min=2,max=100"`
	CreatedTime time.Time     `bson:"created_time" json:"created_time"`
	UpdatedTime time.Time     `bson:"updated_time" json:"updated_time"`
	CreatedBy   string        `bson:"created_by" json:"created_by"`
	UpdatedBy   string        `bson:"updated_by" json:"updated_by"`
}
