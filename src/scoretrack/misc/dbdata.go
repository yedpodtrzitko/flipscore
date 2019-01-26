package misc

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/henderjon/jwt"
	db "upper.io/db.v3"
	"upper.io/db.v3/postgresql"
)

var dbSettings = postgresql.ConnectionURL{
	Database: os.Getenv("DB_NAME"),
	Host:     os.Getenv("DB_HOST"),
	User:     os.Getenv("DB_USER"),
	Password: os.Getenv("DB_PASS"),
}

type ScoreDBRecord struct {
	//Id     uint      `db:"id"`
	GameID    string    `db:"game_id"`
	Score     uint      `db:"score"`
	Player    string    `db:"player"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

type GameKeyRecord struct {
	GameID  string `db:"id"`
	GameKey string `db:"game_key"`
}

var redisDateFormat = "2006-01-02 15:03:04"

func GetDBSession() (db.Database, error) {
	dbSession, err := postgresql.Open(dbSettings)
	dbSession.SetLogging(true)
	return dbSession, err
}

func GetGameKey(gameID string) (string, error) {
	var existingRecord GameKeyRecord
	dbsess, err := GetDBSession()
	defer dbsess.Close()

	keyCollection := dbsess.Collection("game_key")
	err = keyCollection.Find("id", gameID).One(&existingRecord)
	if err != nil {
		return "", errors.New("Game Key not found")
	}

	return existingRecord.GameKey, nil
}

func SaveScore(GameID string, scoreData *jwt.Claims) bool {
	rawScore, err := scoreData.Get("score")
	if err != nil {
		log.Println("failed to get score value")
		return false
	}

	newScore, ok := rawScore.(float64)
	if !ok {
		log.Println("failed to convert score int value")
		return false
	}

	rawGameID, err := scoreData.Get("game_id")
	if GameID != rawGameID {
		log.Println("Game ID value doesnt match")
		return false
	}

	newScoreUint := uint(newScore)
	rawUser, _ := scoreData.Get("player")
	userStr, ok := rawUser.(string)
	if !ok {
		log.Println("failed to get player name")
		return false
	} else {
		//userStr := strings.TrimSpace(userStrRaw)
	}

	rawContent, _ := scoreData.Get("content")
	contentStr, ok := rawContent.(string)

	dbsess, err := GetDBSession()
	defer dbsess.Close()
	if err != nil {
		return false
	}

	var existingRecord ScoreDBRecord
	scoreCollection := dbsess.Collection("score")
	userScoreExists, err := scoreCollection.Find("player", userStr).And("game_id", GameID).Count() //One(&existingRecord)
	if err != nil {
		return false
	}

	if userScoreExists == 0 {
		newRecord := ScoreDBRecord{
			GameID:    GameID,
			Score:     newScoreUint,
			Player:    userStr,
			CreatedAt: time.Now(),
			Content:   contentStr,
		}

		inserted, err := scoreCollection.Insert(newRecord)
		fmt.Println(inserted)
		if err != nil {
			fmt.Println("err on insert")
			return false
		}
		return true
	}

	res := scoreCollection.Find("player", userStr).And("game_id", GameID)
	err = res.One(&existingRecord)
	if err != nil {
		return false
	}

	if existingRecord.Score > newScoreUint {
		return true
	}

	existingRecord.Score = newScoreUint
	existingRecord.Content = contentStr
	err = res.Update(existingRecord)
	if err != nil {
		return false
	}

	return true
}

func GetRedisCache() []ScoreDBRecord {
	// FIX this function
	var scores []ScoreDBRecord

	redisKeys, err := redisClient.Keys("score_*").Result()
	// no redis scores saved
	if err != nil {
		return scores
	}

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
			nativeScore, _ := strconv.ParseUint(redisScore[1].(string), 10, 64)
			score := ScoreDBRecord{
				Player:    redisScore[0].(string),
				Score:     uint(nativeScore),
				CreatedAt: nativeDate,
			}

			scores = append(scores, score)
		}
	}

	return scores
}

func GetScoreList(gameID string) []ScoreDBRecord {
	var scores []ScoreDBRecord
	//scores = GetRedisCache()

	dbsess, err := GetDBSession()
	if err != nil {
		return scores
	}
	defer dbsess.Close()

	res := dbsess.Collection("score").Find("game_id", gameID)
	res = res.OrderBy("score").Limit(5)
	res.All(&scores)

	/*
		if len(scores) > 0 {
			SaveIntoRedis(scores)
		}
		RedisScoreExists(len(scores) > 0)
	*/

	return scores
}

func TestDBConnection() bool {
	_, err := GetDBSession()
	return err == nil
}
