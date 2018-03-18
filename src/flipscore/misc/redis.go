package misc

import (
	"os"
	"github.com/go-redis/redis"
	"strconv"
)

var redisHost = os.Getenv("REDIS_HOST")
var redisDB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))

func InitRedisClient() *redis.Client {
	if redisHost == "" {
		redisHost = "localhost:6379"
	}

	if redisDB == 0 {
		redisDB = 4
	}

	return redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       redisDB,
	})
}

var redisClient = InitRedisClient()
