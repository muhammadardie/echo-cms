package users

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type PublicUsers struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" swaggerignore:"true"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `json:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty"`
}

type Users struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" swaggerignore:"true"`
	Username  string             `json:"username" bson:"username" form:"username" query:"username" swaggerignore:"true"`
	Email     string             `json:"email" bson:"email,omitempty" form:"email" query:"email" validate:"required,email"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty" form:"password" query:"password" validate:"required"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty" swaggerignore:"true"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty" swaggerignore:"true"`
}

type UpdateUser struct {
	Username string `json:"username,omitempty" bson:"username,omitempty" form:"username" query:"username"`
	Email    string `json:"email,omitempty" bson:"email,omitempty" form:"email" query:"email" validate:"omitempty,email"`
	Password string `json:"password,omitempty" bson:"password,omitempty" form:"password" query:"password"` // No "required" validation
}
