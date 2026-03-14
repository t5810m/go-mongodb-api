package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type JobSkill struct {
	ID                       bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	JobID                    bson.ObjectID `bson:"job_id" json:"job_id" validate:"required"`
	SkillID                  bson.ObjectID `bson:"skill_id" json:"skill_id" validate:"required"`
	ProficiencyLevelRequired string        `bson:"proficiency_level_required" json:"proficiency_level_required" validate:"required,oneof=beginner intermediate advanced expert"`
	IsRequired               bool          `bson:"is_required" json:"is_required"`
	CreatedTime              time.Time     `bson:"created_time" json:"created_time"`
	UpdatedTime              time.Time     `bson:"updated_time" json:"updated_time"`
	CreatedBy                string        `bson:"created_by" json:"created_by"`
	UpdatedBy                string        `bson:"updated_by" json:"updated_by"`
}
