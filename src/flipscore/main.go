package main

import (
	"fmt"
	"net/http"
	"./misc"
)

func SaveScoreRoute(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	var signedScore = req.FormValue("score")
	var extractedData = misc.ExtractJWTData(signedScore)

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
	// TODO
	//var scores = misc.GetScoreList()

	fmt.Fprintf(res, "[{}]")
}

func main() {
	http.HandleFunc("/save", SaveScoreRoute)
	http.HandleFunc("/list", GetScoreListRoute)
	http.ListenAndServe(":8080", nil)
}
