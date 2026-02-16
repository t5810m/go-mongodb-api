package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CandidateSkill struct {
	ID               bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CandidateID      bson.ObjectID `bson:"candidate_id" json:"candidate_id" validate:"required"`
	SkillID          bson.ObjectID `bson:"skill_id" json:"skill_id" validate:"required"`
	ProficiencyLevel string        `bson:"proficiency_level" json:"proficiency_level" validate:"required,oneof=beginner intermediate advanced expert"`
	CreatedTime      time.Time     `bson:"created_time" json:"created_time"`
	UpdatedTime      time.Time     `bson:"updated_time" json:"updated_time"`
	CreatedBy        string        `bson:"created_by" json:"created_by"`
	UpdatedBy        string        `bson:"updated_by" json:"updated_by"`
}
