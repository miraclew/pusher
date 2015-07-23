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

var pool *redis.Pool

func SetRedisPool(p *redis.Pool) {
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

func AuthClient(token string) (*Client, error) {
	conn := pool.Get()
	defer conn.Close()

	v, err := redis.StringMap(conn.Do("hgetall", "token:"+token))

	if err != nil {
		return nil, err
	}
	// log.Debug("token:%s %#v", token, v)

	userId, err := strconv.ParseInt(v["user_id"], 10, 64)
	if err != nil {
		return nil, err
	}

	if v["device_type"] == "" {
		v["device_type"] = "2"
	}

	deviceType, err := strconv.ParseInt(v["device_type"], 10, 64)
	if err != nil {
		return nil, err
	}

	client := &Client{}
	client.UserId = userId
	client.Version = v["version"]
	client.DeviceType = int(deviceType)
	// client.NodeId = app.options.nodeId
	return client, nil
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

func RemoveClient(userId int64) {
	conn := pool.Get()
	defer conn.Close()

	k := fmt.Sprintf("client:%d", userId)
	conn.Do("del", k)
}

func (c *Client) Save() error {
	conn := pool.Get()
	defer conn.Close()

	k := fmt.Sprintf("client:%d", c.UserId)
	_, err := conn.Do("hmset", redis.Args{}.Add(k).AddFlat(map[string]string{
		"u":    fmt.Sprintf("%d", c.UserId),
		"dt":   fmt.Sprintf("%d", c.DeviceType),
		"node": fmt.Sprintf("%d", c.NodeId),
		"la":   fmt.Sprintf("%d", c.LastActive),
		"v":    c.Version,
	})...)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Touch(expire int) error {
	conn := pool.Get()
	defer conn.Close()
	k := fmt.Sprintf("client:%d", c.UserId)
	_, err := conn.Do("expire", k, expire)
	return err
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
