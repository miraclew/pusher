package pusher

import (
	"fmt"
)

const (
	DEVICE_TYPE_IOS     = 1
	DEVICE_TYPE_ANDROID = 2
)

type Client struct {
	Token      string
	UserId     int64
	Version    string
	DeviceType int
}

var clients map[int64]*Client

func init() {
	clients = make(map[int64]*Client)
}

func AddClient(client *Client) {
	clients[client.UserId] = client
}

func GetClient(userId int64) *Client {
	return clients[userId]
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

	v, err := VersionCompare(c.Version, gtVersion)
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
