package testimonies

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Testimonies struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username" form:"username" query:"username" validate:"required"`
	Comment   string             `json:"comment" bson:"comment,omitempty" form:"comment" query:"comment" validate:"required"`
	Avatar    string             `json:"avatar,omitempty" bson:"avatar,omitempty" form:"avatar" query:"avatar" validate:"required"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}
