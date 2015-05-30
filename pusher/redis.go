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

func AuthClient(token string) (*Client, error) {
	conn := pool.Get()
	defer conn.Close()

	v, err := redis.StringMap(conn.Do("hgetall", "token:"+token))

	if err != nil {
		return nil, err
	}

	userId, err2 := strconv.ParseInt(v["user_id"], 10, 64)
	if err2 != nil {
		return nil, err
	}

	client := &Client{}
	client.Token = token
	client.UserId = userId
	client.Version = v["version"]
	AddClient(client)

	return client, nil
}
