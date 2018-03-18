package misc

import (
	"upper.io/db.v3/postgresql"
	"time"
	"github.com/henderjon/jwt"
	"upper.io/db.v3"
	"log"
	"os"
	"fmt"
	"strconv"
)

var dbSettings = postgresql.ConnectionURL{
	Database: os.Getenv("DB_NAME"),
	Host:     os.Getenv("DB_HOST"),
	User:     os.Getenv("DB_USER"),
	Password: os.Getenv("DB_PASS"),
}

type ScoreDBRecord struct {
	//Id     uint      `db:"id"`
	Score  uint      `db:"score"`
	Player string    `db:"player"`
	Date   time.Time `db:"created_at"`
}

var redisDateFormat = "2006-01-02 15:03:04"

func SaveScoreDB(scoreData *jwt.Claims) bool {
	rawScore, err := scoreData.Get("score")
	if err != nil {
		log.Println("failed to get score value")
		return false
	}
	scoreInt, ok := rawScore.(float64)
	if !ok {
		log.Println("failed to convert score int value")
		return false
	}

	rawUser, _ := scoreData.Get("player")
	userStr, ok := rawUser.(string)
	if !ok {
		log.Println("failed to get player value")
		return false
	}

	dbItem := ScoreDBRecord{
		Score:  uint(scoreInt),
		Player: userStr,
		Date:   time.Now(),
	}

	dbsess, err := GetDBSession()
	if err != nil {
		return false
	}

	scoreCollection := dbsess.Collection("score")
	scoreCollection.Insert(dbItem)

	defer dbsess.Close()

	return true
}

func ExtractJWTData(jwtSecret string, token string) *jwt.Claims {
	algorithm := jwt.HmacSha256(jwtSecret)

	err := algorithm.Validate(token)
	if err != nil {
		log.Print("not validated")
		panic(err)
	}

	scoreData, err := algorithm.Decode(token)
	if err != nil {
		log.Print("not decoded")

		panic(err)
	}

	return scoreData
}

func GetDBSession() (db.Database, error) {
	dbSession, err := postgresql.Open(dbSettings)
	dbSession.SetLogging(true)
	return dbSession, err
}

func GetScoreList() []ScoreDBRecord {
	var scores []ScoreDBRecord

	redisKeys, err := redisClient.Keys("score_*").Result()
	// no redis scores saved
	if len(redisKeys) == 0 {
		noScores, _ := redisClient.Get("no_score").Result()
		if noScores != "" {
			// no scores anywhere
			return scores
		}
	} else {
		for _, scoreKey := range redisKeys {
			redisScore := redisClient.HMGet(scoreKey, "Player", "Score", "Date").Val()
			nativeDate, _ := time.Parse(redisDateFormat, redisScore[2].(string))
			nativeScore, _ := strconv.ParseUint(redisScore[1].(string),10,64)
			score := ScoreDBRecord{
				Player: redisScore[0].(string),
				Score:  uint(nativeScore),
				Date:   nativeDate,
			}

			scores = append(scores, score)
		}
		return scores
	}

	dbsess, err := GetDBSession()
	if err != nil {
		return scores
	}
	res := dbsess.Collection("score").Find()
	res = res.OrderBy("-score").Limit(5)
	res.All(&scores)
	defer dbsess.Close()

	if len(scores) > 0 {
		SaveIntoRedis(scores)
	}
	RedisScoreExists(len(scores) > 0)

	return scores
}

func RedisScoreExists(hasScore bool) {
	if (hasScore) {
		redisClient.Del("no_score").Result()
	} else {
		redisClient.Set("no_score", "1", 0).Result()
	}
}

func SaveIntoRedis(scores []ScoreDBRecord) {
	for idx, score := range scores {
		redisClient.HSet(fmt.Sprint("score_", idx), "Player", score.Player).Err()
		redisClient.HSet(fmt.Sprint("score_", idx), "Score", score.Score).Err()
		redisClient.HSet(fmt.Sprint("score_", idx), "Date", score.Date.Format(redisDateFormat)).Err()
	}
}

func TestDBConnection() bool {
	_, err := GetDBSession()
	return err == nil
}
