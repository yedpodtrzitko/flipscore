package misc

import (
	"upper.io/db.v3/postgresql"
	"time"
	"github.com/henderjon/jwt"
	"upper.io/db.v3"
	"log"
	"os"
)

var jwtSecret = os.Getenv("JWT_KEY")

var dbSettings = postgresql.ConnectionURL{
	Database: os.Getenv("DB_NAME"),
	Host:     os.Getenv("DB_HOST"),
	User:     os.Getenv("DB_USER"),
	Password: os.Getenv("DB_PASS"),
}

type ScoreDBRecord struct {
	Score  int       `db:"score"`
	Player string    `db:"player"`
	Date   time.Time `db:"created_at"`
}

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
		Score:  int(scoreInt),
		Player: userStr,
		Date:   time.Now(),
	}

	dbsess := GetSession()
	scoreCollection := dbsess.Collection("score")
	scoreCollection.Insert(dbItem)

	defer dbsess.Close()

	return true
}

func ExtractJWTData(token string) *jwt.Claims {
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

func GetSession() db.Database {
	dbSession, err := postgresql.Open(dbSettings)
	dbSession.SetLogging(true)
	if err != nil {
		panic(err)
	}
	return dbSession
}

func GetScoreList() []ScoreDBRecord {

	dbsess := GetSession()

	var scores []ScoreDBRecord

	res := dbsess.Collection("score").Find()
	res = res.OrderBy("-score").Limit(5)

	err := res.All(&scores)
	defer dbsess.Close()

	if err != nil {
		panic(err)
	}

	return scores
}
