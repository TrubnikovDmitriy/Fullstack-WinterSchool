package cache

import (
	"log"
	"github.com/garyburd/redigo/redis"
	"time"
)

var redisPool *redis.Pool

func init() {
	redisPool = initRedisConnectionPool()
}

func initRedisConnectionPool() *redis.Pool {
	redisConnPool := &redis.Pool {
		MaxIdle: 3,
		IdleTimeout: 300 * time.Second,
		MaxActive: 10,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				"localhost:6379",
				redis.DialDatabase(0),
			)
			if err != nil {
				log.Fatal(err)
			}
			// password Auth
			return c, nil
		},
	}

	return redisConnPool
}