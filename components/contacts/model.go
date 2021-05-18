package contacts

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Contacts struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Address   string             `json:"address" bson:"address" form:"address" query:"address" validate:"required"`
	Phone     string             `json:"phone" bson:"phone,omitempty" form:"phone" query:"phone" validate:"required"`
	Mail      string             `json:"mail" bson:"mail,omitempty" form:"mail" query:"mail" validate:"required"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}
