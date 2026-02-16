package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Company struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string        `bson:"name" json:"name" validate:"required,min=3,max=255"`
	Description string        `bson:"description" json:"description" validate:"min=10"`
	Website     string        `bson:"website" json:"website" validate:"required,url"`
	Email       string        `bson:"email" json:"email" validate:"required,email"`
	Phone       string        `bson:"phone" json:"phone" validate:"required,min=10"`
	Address     string        `bson:"address" json:"address"`
	City        string        `bson:"city" json:"city" validate:"required"`
	PostalCode  string        `bson:"postal_code" json:"postal_code"`
	Country     string        `bson:"country" json:"country" validate:"required"`
	LogoUrl     string        `bson:"logo_url" json:"logo_url" validate:"omitempty,url"`
	Verified    bool          `bson:"verified" json:"verified"`
	Active      bool          `bson:"active" json:"active"`
	CreatedTime time.Time     `bson:"created_time" json:"created_time"`
	UpdatedTime time.Time     `bson:"updated_time" json:"updated_time"`
	CreatedBy   string        `bson:"created_by" json:"created_by"`
	UpdatedBy   string        `bson:"updated_by" json:"updated_by"`
}
