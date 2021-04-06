package abouts

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Abouts struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title" form:"title" query:"title"`
	Desc      string             `json:"desc" bson:"desc,omitempty" form:"desc" query:"desc"`
	Image     string             `json:"image,omitempty" bson:"image,omitempty" form:"image" query:"image"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}
