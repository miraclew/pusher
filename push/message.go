package push

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
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

type ServerMsg struct {
}

// server => router
type Message struct {
	Id        int64  `json:"id"`
	Type      int    `json:"type"`
	SubType   int    `db:"sub_type" json:"sub_type"`
	SenderId  int64  `db:"sender_id" json:"sender_id"`
	Receiver  string `json:"receiver"`
	Body      string `json:"body"`
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

func NewMessage(typ int, senderId int64, receiver string, body string, opts string) *Message {
	return &Message{
		Receiver: receiver, Type: typ, Body: body,
		SenderId: senderId, Opts: opts, Timestamp: time.Now().UnixNano() / 1000000,
	}
}

func (m *Message) ParseOpts() *MsgSendOpts {
	return nil
}

func (m *Message) ParseReceivers() ([]int64, error) {
	var receivers []int64
	var errs []string
	rs := strings.Split(m.Receiver, ",")
	for _, r := range rs {
		uid, err := strconv.ParseInt(r, 10, 64)
		if err != nil {
			errs = append(errs, r)
			continue
		}
		receivers = append(receivers, uid)
	}

	var err error
	if len(errs) > 0 {
		err = errors.New(fmt.Sprintf("receivers malform: %s", strings.Join(errs, ",")))
	}

	return receivers, err
}

func (m *Message) Save() error {
	res, err := db.NamedExec(`INSERT INTO messages (type, sender_id, receiver, body, opts, timestamp) VALUES (:type, :sender_id, :receiver, :body, :opts, :timestamp)`, m)
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
