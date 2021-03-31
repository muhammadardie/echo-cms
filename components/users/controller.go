package users

import (
	"context"
	DB "github.com/muhammadardie/echo-cms/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collName string = "users"

func GetUserId(email string) (primitive.ObjectID, error) {
	ctx := context.Background()

	db, err := DB.Connect()
	if err != nil {
		return primitive.NilObjectID, err
	}

	selector := bson.M{"email": email}

	var record Users

	if err = db.Collection(collName).FindOne(ctx, selector).Decode(&record); err != nil {
		return primitive.NilObjectID, err
	}

	return record.ID, nil
}
