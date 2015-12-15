package push

import (
	"encoding/json"
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

	MSG_OPT_OFFLINE_DISABLE = 0
	MSG_OPT_OFFLINE_DEFAULT = 1

	MSG_OPT_ACK_DISABLE = 0
	MSG_OPT_ACK_DEFAULT = 1

	MSG_OPT_APN_DISABLE     = 0
	MSG_OPT_APN_DEFAULT     = 1
	MSG_OPT_APN_NOTIFY_ONLY = 2
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
	Format    int    `json:"format"`
	ChatId    int64  `db:"chat_id" json:"chat_id"`
	Body      string `json:"body"`
	Opts      string `json:"opts"`
	Extra     string `json:"extra"`
	Timestamp int64  `json:"timestamp"` // milseconds
}

type ClientMessage struct {
	Type      int    `json:"type"`
	AckMsgId  string `json:"ack_msg_id"`
	Timestamp int64  `json:"timestamp"`
}

type MsgSendOpts struct {
	TTL         int    `json:"ttl"`
	Alert       string `json:"alert"`
	OfflineFlag int    `json:"offline_flag"`
	AckFlag     int    `json:"ack_flag"`
	ApnFlag     int    `json:"apn_flag"`
	DeviceType  int    `json:"device_type"`
}

func NewMessage(typ int, senderId int64, receiver string, chatId int64, body string, opts string) *Message {
	return &Message{
		Receiver: receiver, Type: typ, Body: body, ChatId: chatId,
		SenderId: senderId, Opts: opts, Timestamp: time.Now().UnixNano() / 1000000,
	}
}

func (m *Message) ParseOpts() *MsgSendOpts {
	opt := &MsgSendOpts{}
	json.Unmarshal([]byte(m.Opts), opt)
	return opt
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
	res, err := db.NamedExec(`INSERT INTO messages (type, sender_id, receiver, chat_id, body, opts, extra, timestamp) VALUES (:type, :sender_id, :receiver, :chat_id, :body, :opts, :extra, :timestamp)`, m)
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

func (m *Message) GetPayload() ([]byte, error) {
	// body := map[string]interface{}{}
	// extra := map[string]interface{}{}
	// err := json.Unmarshal([]byte(m.Body), &body)
	// if err != nil {
	// 	return nil, errors.New("body should be a valid json: " + err.Error())
	// }
	// err := json.Unmarshal([]byte(m.Extra), &extra)
	// if err != nil {
	// 	return nil, errors.New("extra should be a valid json: " + err.Error())
	// }
	return json.Marshal(map[string]interface{}{
		"id":        fmt.Sprintf("%d", m.Id),
		"type":      m.Type,
		"sub_type":  m.SubType,
		"format":    m.Format,
		"chat_id":   m.ChatId,
		"sender_id": fmt.Sprintf("%d", m.SenderId),
		"ttl":       m.ParseOpts().TTL,
		"sent_at":   m.Timestamp / 1000,
		"body":      m.Body,
		"extra":     m.Extra,
	})
}

func FindMessage(id int64) (*Message, error) {
	msg := &Message{}
	err := db.Get(msg, "SELECT * FROM messages WHERE id=?", id)
	return msg, err
}
