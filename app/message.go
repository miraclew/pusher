package app

const (
	MSG_TYPE_DIRECT  = 1
	MSG_TYPE_CHANNEL = 2
	MSG_TYPE_BULK    = 3
)

type Message struct {
	Id        string       `json:"id"`
	Type      int          `json:"type"`
	SenderId  string       `json:"sender_id"`
	ChannelId string       `json:"channel_id"`
	Payload   interface{}  `json:"payload"`
	Opts      *MsgSendOpts `json:"opts"`
	Timestamp int64        `json:"timestamp"`
}

type ClientMessage struct {
	Type      int    `json:"type"`
	AckMsgId  string `json:"ack_msg_id"`
	Timestamp int64  `json:"timestamp"`
}

type MsgSendOpts struct {
	TTL           int    `json:"ttl"`
	Alert         string `json:"alert"`
	OfflineEnable bool   `json:"offline_enable"`
	AckEnable     bool   `json:"ack_enable"`
	ApnEnable     bool   `json:"apn_enable"`
}

func NewMessage(channelId string, typ int, payload interface{}, senderId string, opts *MsgSendOpts) *Message {
	return &Message{
		ChannelId: channelId, Type: typ, Payload: payload,
		SenderId: senderId, Opts: opts, CreatedAt: time.Now(),
	}
}
