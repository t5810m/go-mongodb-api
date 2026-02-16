package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Job struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string        `bson:"title" json:"title" validate:"required,min=5,max=255"`
	Description string        `bson:"description" json:"description" validate:"required,min=20"`
	RecruiterID bson.ObjectID `bson:"recruiter_id" json:"recruiter_id" validate:"required"`
	CompanyID   bson.ObjectID `bson:"company_id" json:"company_id" validate:"required"`
	CategoryID  bson.ObjectID `bson:"category_id" json:"category_id" validate:"required"`
	Location    string        `bson:"location" json:"location" validate:"required,min=3"`
	JobType     string        `bson:"job_type" json:"job_type" validate:"required,oneof=full-time part-time contract freelance"`
	SalaryMin   int           `bson:"salary_min" json:"salary_min" validate:"required,gt=0"`
	SalaryMax   int           `bson:"salary_max" json:"salary_max" validate:"required,gt=0"`
	Status      string        `bson:"status" json:"status" validate:"required,oneof=active closed draft"`
	Active      bool          `bson:"active" json:"active"`
	CreatedTime time.Time     `bson:"created_time" json:"created_time"`
	UpdatedTime time.Time     `bson:"updated_time" json:"updated_time"`
	CreatedBy   string        `bson:"created_by" json:"created_by"`
	UpdatedBy   string        `bson:"updated_by" json:"updated_by"`
}
