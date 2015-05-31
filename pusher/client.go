package pusher

type Client struct {
	Token   string
	UserId  int64
	Version string
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

	v, err := VersionCompare(c.Version, "2.0.5")
	if err != nil {
		return false
	}

	return v > 0
}
