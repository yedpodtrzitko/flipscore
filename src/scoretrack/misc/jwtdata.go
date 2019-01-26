package misc

import (
	"log"

	"github.com/henderjon/jwt"
)

func Extract(jwtKey string, jwtData string) *jwt.Claims {
	algorithm := jwt.HmacSha256(jwtKey)

	err := algorithm.Validate(jwtData)
	if err != nil {
		log.Print("jwt not valid")
		panic(err)
	}

	scoreData, err := algorithm.Decode(jwtData)
	if err != nil {
		log.Print("jwt not decoded")
		panic(err)
	}

	return scoreData
}
