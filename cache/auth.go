package cache

import (
	"../models"
	"../services"
	"github.com/satori/go.uuid"
	"github.com/garyburd/redigo/redis"
	"log"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"time"
)


func CreateCode(oauth *models.OAuth) *uuid.UUID {

	deleteTokens(oauth.PersonID)

	conn := redisPool.Get()
	defer conn.Close()

	code, _ := uuid.NewV4()

	jsonAuth, err := json.Marshal(oauth)
	if err != nil {
		log.Print(err)
		return nil
	}

	_, err = redis.String(conn.Do(
		"HMSET", "code:" + code.String(),
		"access", string(jsonAuth),
		"person_id", oauth.PersonID.String(),
	))
	if err != nil {
		log.Print(err)
		return nil
	}
	redis.String(conn.Do("EXPIRE", "code:" + code.String(), 5 * time.Minute))

	return &code
}

func ActivateCode(code *uuid.UUID) (newAccess *string, newRefresh *string) {

	conn := redisPool.Get()
	defer conn.Close()

	personUUID, err := redis.String(conn.Do("HGET", "code:" + code.String(), "person_id"))
	if err != nil {
		// Такого кода не существует
		return nil, nil
	}

	access, err := redis.String(conn.Do("HGET", "code:" + code.String(), "access"))
	if err != nil {
		// Возможно кто-то пытается воспользоваться кодом второй раз,
		// на всякий случай, лучше инвалидировать токены
		log.Printf("Access to the token more than once (personID: %s)\n", personUUID)
		deleteTokens(uuid.FromStringOrNil(personUUID))
		return nil, nil
	}

	// Чтобы предотвратить повторное чтение удаляем access-token из code
	redis.String(conn.Do("HDEL", "code:" + code.String(), "access"))

	// Формируем две новые записи: (refresh:uuid, models.OAuth), (person_id:uuid, refresh_uuid)
	newRefresh = generateRefreshToken()
	redis.String(conn.Do("SET", "person_id:" + personUUID, *newRefresh))
	redis.String(conn.Do("SET", "refresh:" + *newRefresh, access))

	// На основе модели генерируем новый access token
	auth := unmarshalAuth(access)
	return generateAccessToken(auth), newRefresh
}

func RefreshToken(refreshToken string) (newAccess *string, newRefresh *string) {

	conn := redisPool.Get()
	defer conn.Close()

	// Получить AccessToken и сразу удалить запись
	access, err := redis.String(conn.Do("GET", "refresh:" + refreshToken))
	if err != nil {
		return nil, nil
	}
	redis.String(conn.Do("DEL", "refresh:" + refreshToken))

	// Разобарть Access Token
	oauth := unmarshalAuth(access)
	newRefresh = generateRefreshToken()

	// Обновить и создать записи
	redis.String(conn.Do("SET", "person_id:" + oauth.PersonID.String(), *newRefresh))
	redis.String(conn.Do("SET", "refresh:" + *newRefresh, access))

	return generateAccessToken(oauth), newRefresh
}


func deleteTokens(personID uuid.UUID) {

	conn := redisPool.Get()
	defer conn.Close()

	refresh, err := redis.String(conn.Do("GET", "person_id:" + personID.String()))
	if err == nil {
		redis.String(conn.Do("DEL", "person_id:" + personID.String()))
		redis.String(conn.Do("DEL", "refresh:" + refresh))
	}
}

func generateRefreshToken() *string {
	refresh, _ := uuid.NewV4()
	stringValue := refresh.String()
	return &stringValue
}

func generateAccessToken(oauth *models.OAuth) *string {

	tokenAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(serv.GetConfig().AccessTokenTTL).Unix(),

		"person_id":  oauth.PersonID,
		"first_name": oauth.FirstName,
		"last_name":  oauth.LastName,
		"staff":     oauth.Staff,
		"app_id":	 oauth.AppID,
		"scope":	 oauth.Scope,
	})

	tokenString, err := tokenAccess.SignedString([]byte(serv.GetConfig().SignKey))
	if err != nil {
		log.Print(err)
		return nil
	}
	return &tokenString
}

func unmarshalAuth(authJSON string) *models.OAuth {

	auth := models.OAuth{}
	json.Unmarshal([]byte(authJSON), &auth)
	return &auth
}
