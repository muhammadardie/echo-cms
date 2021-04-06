package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Users struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" swaggerignore:"true"`
	Username  string             `json:"username" bson:"username" form:"username" query:"username" swaggerignore:"true"`
	Email     string             `json:"email" bson:"email,omitempty" form:"email" query:"email" validate:"required"`
	Password  string             `json:"password" bson:"password,omitempty" form:"password" query:"password" validate:"required"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty" swaggerignore:"true"`
}
