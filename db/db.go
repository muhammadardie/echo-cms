package db

import (
	"context"
	_ "github.com/joho/godotenv/autoload" // read .env on import
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"sync"
)

/* Used to create a singleton object of MongoDB client. */
var clientDatabase *mongo.Database
var clientInstanceError error

//Used to execute client creation procedure only once.
var mongoOnce sync.Once

func Connect() (*mongo.Database, error) {
	mongoUrl := os.Getenv("MONGODB_URL")
	mongoDBName := os.Getenv("MONGODB_NAME")
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(mongoUrl)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientDatabase = client.Database(mongoDBName)
	})

	return clientDatabase, clientInstanceError
}
