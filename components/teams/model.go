package teams

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Teams struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" form:"name" query:"name" validate:"required"`
	Position  string             `json:"position" bson:"position,omitempty" form:"position" query:"position" validate:"required"`
	Image     string             `json:"image,omitempty" bson:"image,omitempty" form:"image" query:"image" validate:"required"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}
