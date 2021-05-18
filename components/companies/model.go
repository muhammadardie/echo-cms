package companies

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Companies struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title" form:"title" query:"title" validate:"required"`
	Desc   	  string             `json:"desc" bson:"desc,omitempty" form:"desc" query:"desc" validate:"required"`
	Image     string             `json:"image,omitempty" bson:"image,omitempty" form:"image" query:"image" validate:"required"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}
