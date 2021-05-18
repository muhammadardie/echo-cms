package socmeds

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Socmeds struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" form:"name" query:"name" validate:"required"`
	Icon      string             `json:"icon" bson:"icon,omitempty" form:"icon" query:"icon" validate:"required"`
	Url       string             `json:"url" bson:"url,omitempty" form:"url" query:"url" validate:"required"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}
