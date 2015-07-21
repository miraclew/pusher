package push

import (
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	ROUTE_TYPE_DIRECT  = 1
	ROUTE_TYPE_CHANNEL = 2
	ROUTE_TYPE_BULK    = 3

	MSG_TYPE_ACK = 6001
)

var db *sqlx.DB

func SetDb(d *sqlx.DB) {
	db = d
}

type Message struct {
	Id        string `json:"id"`
	Type      int    `json:"type"`
	SenderId  int64  `json:"sender_id"`
	Receiver  string `json:"receiver"`
	Payload   string `json:"payload"`
	Opts      string `json:"opts"`
	Timestamp int64  `json:"timestamp"`
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

func NewMessage(typ int, senderId int64, receiver string, payload string, opts string) *Message {
	return &Message{
		Receiver: receiver, Type: typ, Payload: payload,
		SenderId: senderId, Opts: opts, Timestamp: time.Now().UnixNano(),
	}
}

func (m *Message) ParseOpts() *MsgSendOpts {
	return nil
}

func (m *Message) Save() error {
	_, err := db.NamedExec(`INSERT INTO messages (type, sender_id, receiver, payload, opts, timestamp) VALUES (:type, :sender_id, :receiver, :payload, :opts, :timestamp)`, m)
	return err
}

func FindMessage(id string) (*Message, error) {
	msg := &Message{}
	err := db.Get(msg, "SELECT * FROM messages WHERE id=$1", id)
	return msg, err
}
