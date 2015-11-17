package push

import (
	"github.com/garyburd/redigo/redis"
	"testing"
	"time"
)

func setupRedis() {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "ubuntu:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	SetRedisPool(pool)
}

func _TestClient(t *testing.T) {
	setupRedis()
	client := &Client{}
	client.UserId = 111
	client.Version = "2.0"
	client.DeviceType = 1
	client.NodeId = 2

	err := client.Save()
	if err != nil {
		t.Fail()
	}
}
