package main

import (
	"fmt"
	"net/http"
	"./misc"
	"encoding/json"
	"os"
	"log"
)

var jwtSecret = os.Getenv("JWT_KEY")


func SaveScoreRoute(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	signedScore := req.FormValue("score")
	extractedData := misc.ExtractJWTData(jwtSecret, signedScore)

	if extractedData == nil {
		fmt.Fprintf(res, "signature check failed")
		res.WriteHeader(400)
		return
	}

	if !misc.SaveScoreDB(extractedData) {
		fmt.Fprintf(res, "failed to save score")
		res.WriteHeader(500)
		return
	}

	fmt.Fprintf(res, "score saved %s", req.URL.Path)
}

func GetScoreListRoute(res http.ResponseWriter, req *http.Request) {
	scores := misc.GetScoreList()
	scoreBytes, _ := json.Marshal(&scores)
	scoreJson := string(scoreBytes[:])

	fmt.Fprintf(res, scoreJson)
}

func main() {
	if jwtSecret == "" {
		log.Print("JWT_KEY not defined")
		os.Exit(1)
	}

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
	http.ListenAndServe(fmt.Sprint(":", serverPort), nil)
}
