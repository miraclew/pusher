package push

import (
	"coding.net/miraclew/pusher/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

const (
	DEVICE_TYPE_IOS     = 1
	DEVICE_TYPE_ANDROID = 2
	ONLINE_TIMEOUT      = 30
)

var pool redis.Pool

func SetRedisPool(p redis.Pool) {
	pool = p
}

type Client struct {
	UserId     int64
	DeviceType int
	Version    string
	// Token      string
	NodeId     int
	LastActive int64
}

func GetClient(userId int64) (*Client, error) {
	conn := pool.Get()
	defer conn.Close()

	sm, err := redis.StringMap(conn.Do("hgetall", fmt.Sprintf("client:%d", userId)))
	if err != nil {
		return nil, err
	}

	deviceType, _ := strconv.ParseInt(sm["dt"], 10, 64)
	nodeId, _ := strconv.ParseInt(sm["node"], 10, 64)
	lastActive, _ := strconv.ParseInt(sm["la"], 10, 64)

	return &Client{
		UserId:     userId,
		DeviceType: int(deviceType),
		Version:    sm["v"],
		// Token:      sm["t"],
		LastActive: lastActive,
		NodeId:     int(nodeId),
	}, nil
}

func (c *Client) Save() error {
	conn := pool.Get()
	defer conn.Close()

	args := redis.Args{}
	args.Add(fmt.Sprintf("client:%d", c.UserId))
	args.AddFlat(map[string]string{
		"u":    fmt.Sprintf("%d", c.UserId),
		"dt":   fmt.Sprintf("%d", c.DeviceType),
		"node": fmt.Sprintf("%d", c.NodeId),
		"la":   fmt.Sprintf("%d", c.LastActive),
		"v":    c.Version,
	})

	_, err := conn.Do("hmset", args...)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) IsOnline() bool {
	return c.NodeId > 0 && (time.Now().Unix()-c.LastActive > ONLINE_TIMEOUT)
}

func (c *Client) SupportAck() bool {
	if len(c.Version) <= 0 {
		return false
	}

	gtVersion := "2.0.5"
	if c.DeviceType == DEVICE_TYPE_IOS {
		gtVersion = "2.0.0"
	} else if c.DeviceType == DEVICE_TYPE_ANDROID {

	}

	v, err := util.VersionCompare(c.Version, gtVersion)
	if err != nil {
		return false
	}

	return v > 0
}

func (c *Client) DeviceTypeName() string {
	if c.DeviceType == DEVICE_TYPE_IOS {
		return "iOS"
	} else if c.DeviceType == DEVICE_TYPE_ANDROID {
		return "Android"
	} else {
		return fmt.Sprintf("Unknow: %d", c.DeviceType)
	}
}
