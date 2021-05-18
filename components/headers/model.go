package headers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Headers struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Image     string             `json:"image,omitempty" bson:"image,omitempty" form:"image" query:"image" validate:"required"`
	Page      string             `json:"page" bson:"page" form:"page" query:"page" validate:"required"`
	Tagline   string             `json:"tagline" bson:"tagline,omitempty" form:"tagline" query:"content" validate:"required"`
	Tagdesc   string             `json:"tagdesc,omitempty" bson:"tagdesc,omitempty" form:"tagdesc" query:"image" validate:"required"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}
