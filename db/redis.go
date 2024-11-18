package db

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func InitRedis() *redis.Client {
	//Initializing redis
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)

	return rdb
}
