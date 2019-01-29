package misc

import (
	"log"

	"github.com/henderjon/jwt"
)

func ExtractJWT(jwtKey string, jwtData string) *jwt.Claims {
	algorithm := jwt.HmacSha256(jwtKey)
	err := algorithm.Validate(jwtData)
	if err != nil {
		log.Println(err)
		log.Println("jwt not valid")
		return nil
	}

	scoreData, err := algorithm.Decode(jwtData)
	if err != nil {
		log.Println(err)
		log.Println("jwt not decoded")
		return nil
	}

	return scoreData
}
