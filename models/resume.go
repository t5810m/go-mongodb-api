package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Resume struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CandidateID  bson.ObjectID `bson:"candidate_id" json:"candidate_id" validate:"required"`
	FileUrl      string        `bson:"file_url" json:"file_url" validate:"required,url"`
	FileName     string        `bson:"file_name" json:"file_name" validate:"required,min=3"`
	UploadedTime time.Time     `bson:"uploaded_time" json:"uploaded_time"`
	UpdatedTime  time.Time     `bson:"updated_time" json:"updated_time"`
	CreatedBy    string        `bson:"created_by" json:"created_by"`
	UpdatedBy    string        `bson:"updated_by" json:"updated_by"`
}
