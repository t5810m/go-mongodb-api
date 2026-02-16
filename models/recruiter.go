package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Recruiter struct {
	ID          bson.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName   string         `bson:"first_name" json:"first_name" validate:"required,min=2,max=100"`
	LastName    string         `bson:"last_name" json:"last_name" validate:"required,min=2,max=100"`
	Email       string         `bson:"email" json:"email" validate:"required,email"`
	Password    string         `bson:"password" json:"password" validate:"required,min=8"`
	Phone       string         `bson:"phone" json:"phone" validate:"required,min=10"`
	CompanyID   *bson.ObjectID `bson:"company_id" json:"company_id"` // nil = independent recruiter
	Verified    bool           `bson:"verified" json:"verified"`
	Active      bool           `bson:"active" json:"active"`
	CreatedTime time.Time      `bson:"created_time" json:"created_time"`
	UpdatedTime time.Time      `bson:"updated_time" json:"updated_time"`
	CreatedBy   string         `bson:"created_by" json:"created_by"`
	UpdatedBy   string         `bson:"updated_by" json:"updated_by"`
}
