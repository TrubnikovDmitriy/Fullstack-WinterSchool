package handlers

import (
	"../services"
	"github.com/valyala/fasthttp"
	"github.com/satori/go.uuid"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
	"strconv"
)

const CookieAccess = "ws_auth"
const CookieRefresh = "ws_auth_refresh"


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

func getIntFromBytes(number []byte, defaulting int) int {
	integer, err := strconv.Atoi(string(number))
	if err != nil {
		return defaulting
	}
	return integer
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

func GetClaimsFromCookie(ctx *fasthttp.RequestCtx) (claims jwt.MapClaims, err *serv.ErrorCode) {

	accessToken := string(ctx.Request.Header.Cookie(CookieAccess))
	if len(accessToken) == 0 {
		err = serv.NewUnauthorized()
		return
	}

	claimsPtr := ParseAccessToken(&accessToken)
	if claimsPtr != nil {
		err = serv.NewBadRequest("Token is not valid")
		return
	}

	if IsExpire(claimsPtr) {
		err = &serv.ErrorCode{
			Code: fasthttp.StatusForbidden,
			Message: "Cookie is expired",
			Link: serv.GetConfig().Href + "/v1/oauth/refresh",
		}
		return
	}

	claims, err = *claimsPtr, nil
	return
}