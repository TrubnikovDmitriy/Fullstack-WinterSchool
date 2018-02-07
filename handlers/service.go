package handlers

import (
	"../services"
	"github.com/valyala/fasthttp"
	"github.com/satori/go.uuid"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)


func getPathID(strID string) (uuid.UUID, *serv.ErrorCode) {
	id, err := uuid.FromString(strID)
	if err != nil {
		return uuid.Nil, &serv.ErrorCode{
			Code: fasthttp.StatusBadRequest,
			Message: "Incorrect path variable: '" + strID + "'",
			Link: "Ссылка на документацию",
		}
	}
	return id, nil
}

func ParseAccessToken(accessToken *string) *jwt.MapClaims {

	token, err := jwt.Parse(*accessToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, nil
		}
		return []byte(serv.GetConfig().SignKey), nil
	})
	if err != nil {
		log.Printf("Can't parse token")
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Printf("Token is not valid token: %s", token)
		return nil
	}

	return &claims
}

func IsExpire(claims *jwt.MapClaims) bool {
	expire := time.Unix(int64((*claims)["exp"].(float64)), 0)
	return expire.Before(time.Now())
}