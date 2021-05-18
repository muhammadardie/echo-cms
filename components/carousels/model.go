package carousels

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Carousels struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Tagline   string             `json:"tagline" bson:"tagline" form:"tagline" query:"tagline" validate:"required"`
	Tagdesc   string             `json:"tagdesc" bson:"tagdesc,omitempty" form:"tagdesc" query:"content" validate:"required"`
	Image     string             `json:"image,omitempty" bson:"image,omitempty" form:"image" query:"image" validate:"required"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}
