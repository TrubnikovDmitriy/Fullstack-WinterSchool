package cache

import (
	"log"
	"github.com/garyburd/redigo/redis"
	"../services"
	"time"
)


var redisConnPools []*redis.Pool

func init() {
	redisConnPools = initRedisConnectionPool()
}

func initRedisConnectionPool() []*redis.Pool {

	for _, conf := range serv.GetConfig().Redis {
		connPool := &redis.Pool {
			MaxIdle: 3,
			IdleTimeout: 300 * time.Second,
			MaxActive: 50,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial(
					"tcp",
					conf.Address,
					redis.DialDatabase(0),
				)
				if err != nil {
					log.Fatal(err)
				}
				return c, nil
			},
		}
		redisConnPools = append(redisConnPools, connPool)
	}

	return redisConnPools
}

func sharedKeyByString(uuid string) *redis.Pool {
	return redisConnPools[int(uuid[0]) % len(redisConnPools)]
}