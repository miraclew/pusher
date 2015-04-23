package pusher

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

var ()

func newRedisPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			// if _, err := c.Do("AUTH", password); err != nil {
			//     c.Close()
			//     return nil, err
			// }
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func GetUserIdByToken(token string) (int64, error) {
	v, err := redis.StringMap(pool.Get().Do("hgetall", "token:"+token))

	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(v["user_id"], 10, 64)
}
