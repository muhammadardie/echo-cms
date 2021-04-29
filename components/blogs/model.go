package blogs

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Blogs struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title" form:"title" query:"title" validate:"required"`
	Content   string             `json:"content" bson:"content,omitempty" form:"content" query:"content" validate:"required"`
	Image     string             `json:"image,omitempty" bson:"image,omitempty" form:"image" query:"image" validate:"required"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
}
