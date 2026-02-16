package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// User represents an admin user who manages the platform
type User struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string        `bson:"first_name" json:"first_name" validate:"required,min=2,max=100"`
	LastName          string        `bson:"last_name" json:"last_name" validate:"required,min=2,max=100"`
	Email             string        `bson:"email" json:"email" validate:"required,email"`
	Password          string        `bson:"password" json:"password" validate:"required,min=8"`
	Verified          bool          `bson:"verified" json:"verified"`
	Active            bool          `bson:"active" json:"active"`
	TermsAccepted     bool          `bson:"terms_accepted" json:"terms_accepted"`
	LastTermsAccepted time.Time     `bson:"last_terms_accepted" json:"last_terms_accepted"`
	LastLoginTime     time.Time     `bson:"last_login_time" json:"last_login_time"`
	CreatedTime       time.Time     `bson:"created_time" json:"created_time"`
	UpdatedTime       time.Time     `bson:"updated_time" json:"updated_time"`
	CreatedBy         string        `bson:"created_by" json:"created_by"`
	UpdatedBy         string        `bson:"updated_by" json:"updated_by"`
}
