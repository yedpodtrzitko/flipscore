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
	GameID         string `db:"id"`
	GameKey        string `db:"game_key"`
	ScoreAscending bool   `db:"score_ascending"`
	ScoreInterval  uint   `db:"score_interval"`
}

var redisDateFormat = "2006-01-02 15:03:04"

func GetDBSession() db.Database {
	dbSession, err := postgresql.Open(dbSettings)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if os.Getenv("DB_LOGGING") != "" {
		dbSession.SetLogging(true)
		fmt.Println("db logging on")
	}
	return dbSession
}

func GetGameInfo(gameID string) (GameKeyRecord, error) {
	var existingRecord GameKeyRecord
	dbsess := GetDBSession()
	defer dbsess.Close()

	keyCollection := dbsess.Collection("game_key")
	err := keyCollection.Find("id", gameID).One(&existingRecord)
	if err != nil {
		return existingRecord, errors.New("Game Key not found")
	}

	return existingRecord, nil
}

func SaveScore(GameInfo GameKeyRecord, scoreData *jwt.Claims) bool {
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
	if GameInfo.GameID != rawGameID {
		log.Println("Game ID value doesnt match")
		return false
	}

	newScoreUint := uint(newScore)
	rawUser, _ := scoreData.Get("player")
	userStr, ok := rawUser.(string)
	if !ok || len(userStr) < 2 {
		log.Println("failed to get player name")
		return false
	} else {
		//userStr := strings.TrimSpace(userStrRaw)
	}

	rawContent, _ := scoreData.Get("content")
	contentStr, ok := rawContent.(string)

	dbsess := GetDBSession()
	defer dbsess.Close()

	var now = time.Now()
	dateStr := now.Format("2006-01-02")
	scoreCollection := dbsess.Collection("score")
	userScoreExists := scoreCollection.Find("player", userStr).And("game_id", GameInfo.GameID)

	if GameInfo.ScoreInterval == 24 {
		userScoreExists = userScoreExists.And(db.Raw("CAST(created_at as date) = ?", dateStr))
	}

	exists, err := userScoreExists.Count()
	if err != nil {
		return false
	}

	if exists == 0 {
		newRecord := ScoreDBRecord{
			GameID:    GameInfo.GameID,
			Score:     newScoreUint,
			Player:    userStr,
			CreatedAt: now,
			Content:   contentStr,
		}

		_, err := scoreCollection.Insert(newRecord)
		if err != nil {
			fmt.Println("err on insert")
			return false
		}
		return true
	}

	//res := scoreCollection.Find("player", userStr).And("game_id", GameID)
	//err = res.One(&existingRecord)
	var existingRecord ScoreDBRecord
	err = userScoreExists.One(&existingRecord)
	if err != nil {
		return false
	}

	if GameInfo.ScoreAscending && existingRecord.Score >= newScoreUint {
		return true
	} else if !GameInfo.ScoreAscending && existingRecord.Score <= newScoreUint {
		return true
	}

	existingRecord.Score = newScoreUint
	existingRecord.Content = contentStr
	existingRecord.CreatedAt = now
	err = userScoreExists.Update(existingRecord)
	return err != nil
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

	dbsess := GetDBSession()
	defer dbsess.Close()

	gameInfo, err := GetGameInfo(gameID)
	if err != nil {
		return scores
	}

	orderCond := "-score"
	if gameInfo.ScoreAscending {
		orderCond = "score"
	}
	res := dbsess.Collection("score").Find("game_id", gameID)
	res = res.OrderBy(orderCond).Limit(5)
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
	conn := GetDBSession()
	defer conn.Close()
	return conn != nil
}
