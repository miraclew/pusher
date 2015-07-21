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
	Id        int64  `json:"id"`
	Type      int    `json:"type"`
	SenderId  int64  `db:"sender_id" json:"sender_id"`
	Receiver  string `json:"receiver"`
	Payload   string `json:"payload"`
	Opts      string `json:"opts"`
	Timestamp int64  `json:"timestamp"` // milseconds
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
		SenderId: senderId, Opts: opts, Timestamp: time.Now().UnixNano() / 1000000,
	}
}

func (m *Message) ParseOpts() *MsgSendOpts {
	return nil
}

func (m *Message) Save() error {
	res, err := db.NamedExec(`INSERT INTO messages (type, sender_id, receiver, payload, opts, timestamp) VALUES (:type, :sender_id, :receiver, :payload, :opts, :timestamp)`, m)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	m.Id = id
	return nil
}

func FindMessage(id int64) (*Message, error) {
	msg := &Message{}
	err := db.Get(msg, "SELECT * FROM messages WHERE id=?", id)
	return msg, err
}
