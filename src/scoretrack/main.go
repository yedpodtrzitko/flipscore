package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"./misc"
)

func SaveScoreRoute(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	gameID := req.FormValue("game_id")
	//fmt.Println(gameID)
	gameInfo, err := misc.GetGameInfo(gameID)
	if err != nil {
		fmt.Fprintf(res, "game not found")
		res.WriteHeader(500)
		return
	}

	signedData := req.FormValue("jwt_data")
	extractedData := misc.ExtractJWT(gameInfo.GameKey, signedData)
	if extractedData == nil {
		fmt.Fprintf(res, "jwt signature check failed")
		res.WriteHeader(400)
		return
	}

	if !misc.SaveScore(gameInfo, extractedData) {
		fmt.Fprintf(res, "failed to save score")
		res.WriteHeader(500)
		return
	}

	//misc.SaveIntoRedis()
	// reset redis no_score check
	//misc.RedisScoreExists(true)

	res.WriteHeader(201)
	fmt.Fprintf(res, "score saved %s", req.URL.Path)
}

func GetScoreListRoute(res http.ResponseWriter, req *http.Request) {
	keys, ok := req.URL.Query()["gameID"]
	if !ok || len(keys[0]) < 1 {
		fmt.Fprintf(res, "missing gameID parameter")
		res.WriteHeader(400)
		return
	}

	scores := misc.GetScoreList(keys[0])
	scoreBytes, _ := json.Marshal(&scores)
	scoreJson := string(scoreBytes[:])

	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, scoreJson)
}

func GetIndex(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "ok")
}

func main() {
	serverPort := os.Getenv("SERVER_PORT")
	if len(serverPort) == 0 {
		serverPort = "4980"
	}

	if !misc.TestDBConnection() {
		log.Print("wrong DB configuration")
		os.Exit(1)
	}

	log.Println("serving on port", serverPort)

	http.HandleFunc("/save", SaveScoreRoute)
	http.HandleFunc("/list", GetScoreListRoute)
	http.HandleFunc("/", GetIndex)
	http.ListenAndServe(fmt.Sprint(":", serverPort), nil)
}
