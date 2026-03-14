package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// User represents a platform user with a specific role (admin, candidate, recruiter)
type User struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string        `bson:"first_name" json:"first_name" validate:"required,min=2,max=100"`
	LastName          string        `bson:"last_name" json:"last_name" validate:"required,min=2,max=100"`
	Email             string        `bson:"email" json:"email" validate:"required,email"`
	Password          string        `bson:"password" json:"password,omitempty" validate:"required,min=8"`
	Phone             string        `bson:"phone,omitempty" json:"phone,omitempty"`
	Role              string        `bson:"role" json:"role" validate:"required,oneof=admin candidate recruiter"`
	CompanyName       string        `bson:"company_name,omitempty" json:"company_name,omitempty"`
	Verified          bool          `bson:"verified" json:"verified"`
	Active            bool          `bson:"active" json:"active"`
	TermsAccepted     bool          `bson:"terms_accepted" json:"terms_accepted"`
	LastTermsAccepted *time.Time    `bson:"last_terms_accepted,omitempty" json:"last_terms_accepted,omitempty"`
	LastLoginTime     *time.Time    `bson:"last_login_time,omitempty" json:"last_login_time,omitempty"`
	CreatedTime       time.Time     `bson:"created_time" json:"created_time"`
	UpdatedTime       time.Time     `bson:"updated_time" json:"updated_time"`
	CreatedBy         string        `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedBy         string        `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
}

type UserResponse struct {
	ID                bson.ObjectID `json:"id,omitempty"`
	FirstName         string        `json:"first_name"`
	LastName          string        `json:"last_name"`
	Email             string        `json:"email"`
	Phone             string        `json:"phone,omitempty"`
	Role              string        `json:"role"`
	CompanyName       string        `json:"company_name,omitempty"`
	Verified          bool          `json:"verified"`
	Active            bool          `json:"active"`
	TermsAccepted     bool          `json:"terms_accepted"`
	LastTermsAccepted *time.Time    `json:"last_terms_accepted,omitempty"`
	LastLoginTime     *time.Time    `json:"last_login_time,omitempty"`
	CreatedTime       time.Time     `json:"created_time"`
	UpdatedTime       time.Time     `json:"updated_time"`
	CreatedBy         string        `json:"created_by,omitempty"`
	UpdatedBy         string        `json:"updated_by,omitempty"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:                u.ID,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		Email:             u.Email,
		Phone:             u.Phone,
		Role:              u.Role,
		CompanyName:       u.CompanyName,
		Verified:          u.Verified,
		Active:            u.Active,
		TermsAccepted:     u.TermsAccepted,
		LastTermsAccepted: u.LastTermsAccepted,
		LastLoginTime:     u.LastLoginTime,
		CreatedTime:       u.CreatedTime,
		UpdatedTime:       u.UpdatedTime,
		CreatedBy:         u.CreatedBy,
		UpdatedBy:         u.UpdatedBy,
	}
}
