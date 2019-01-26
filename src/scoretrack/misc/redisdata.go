package misc

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis"
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

func RedisScoreExists(hasScore bool) {
	if hasScore {
		redisClient.Del("no_score").Result()
	} else {
		redisClient.Set("no_score", "1", 0).Result()
	}
}

func SaveIntoRedis(scores []ScoreDBRecord) {
	for idx, score := range scores {
		redisClient.HSet(fmt.Sprint("score_", idx), "Player", score.Player).Err()
		redisClient.HSet(fmt.Sprint("score_", idx), "Score", score.Score).Err()
		redisClient.HSet(fmt.Sprint("score_", idx), "CreatedAt", score.CreatedAt.Format(redisDateFormat)).Err()
	}
}
