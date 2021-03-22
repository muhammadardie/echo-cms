package abouts

import (
	"time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Abouts struct {
    ID        primitive.ObjectID `json:"_id" bson:"_id"`
    Title     string             `json:"title" bson:"title" form:"title" query:"title"`
    Desc      string             `json:"desc" bson:"desc" form:"desc" query:"desc"`
    Image     string             `json:"image" bson:"image" form:"image" query:"image"`
    CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
    UpdatedAt time.Time          `json:"updtedAt" bson:"updated_at"`
}