package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobCategory struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string        `bson:"name" json:"name" validate:"required,min=3,max=100"`
	Description string        `bson:"description" json:"description" validate:"required,min=10,max=500"`
	CreatedTime time.Time     `bson:"created_time" json:"created_time"`
	UpdatedTime time.Time     `bson:"updated_time" json:"updated_time"`
	CreatedBy   string        `bson:"created_by" json:"created_by"`
	UpdatedBy   string        `bson:"updated_by" json:"updated_by"`
}
