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
				redis.DialDatabase(1),
			)
			if err != nil {
				log.Fatal(err)
			}
			// password Auth
			return c, nil
		},
	}

	//conn := redisConnPool.Get()
	//defer conn.Close()
	//
	//str, err := redis.String(conn.Do("PING"))
	//if err != nil {
	//	fmt.Print("Bad")
	//}
	//fmt.Printf("Good: %s\n", str)
	//str, err = redis.String(conn.Do("SET", "test", "world"))
	//if err != nil {
	//	fmt.Print(err)
	//} else {
	//	fmt.Print(str)
	//}
	//conn.Do("APPEND", "test", "Hello")
	//str, err = redis.String(conn.Do("GET", "test"))
	//fmt.Printf("\nGood: %s", str)

	return redisConnPool
}
