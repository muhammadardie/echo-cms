package db

import (
	"github.com/go-redis/redis/v8"
	"os"
)

var client *redis.Client

func InitRedis() *redis.Client {
	//Initializing redis
	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		panic(err)
	}

	return client
}
