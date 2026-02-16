package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Application struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	JobID         bson.ObjectID `bson:"job_id" json:"job_id" validate:"required"`
	CandidateID   bson.ObjectID `bson:"candidate_id" json:"candidate_id" validate:"required"`
	Status        string        `bson:"status" json:"status" validate:"required,oneof=applied under_review rejected accepted withdrawn"`
	RecruiterNote string        `bson:"recruiter_note" json:"recruiter_note"`
	AppliedTime   time.Time     `bson:"applied_time" json:"applied_time"`
	UpdatedTime   time.Time     `bson:"updated_time" json:"updated_time"`
	CreatedBy     string        `bson:"created_by" json:"created_by"`
	UpdatedBy     string        `bson:"updated_by" json:"updated_by"`
}
