package pusher

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	r "github.com/dancannon/gorethink"
	"github.com/garyburd/redigo/redis"
	"io"
	"log"
	"strings"
	"time"
)

const (
	CHANNEL_TYPE_NORMAL = 1
	MESSAGE_TYPE_NORMAL = 1

	MSG_TYPE_ACK = 6001
)

type Channel struct {
	Id        string    `gorethink:"id,omitempty" json:"id"`
	Hash      string    `gorethink:"hash" json:"hash"`
	Members   []string  `gorethink:"members" json:"members"`
	Type      int       `gorethink:"type" json:"type"`
	CreatedAt time.Time `gorethink:"created_at" json:"created_at"`
}

type Message struct {
	Id        string                 `gorethink:"id,omitempty" json:"id"`
	ChannelId string                 `gorethink:"channel_id" json:"channel_id"`
	Type      int                    `gorethink:"type" json:"type"`
	Payload   interface{}            `gorethink:"payload" json:"payload"`
	SenderId  string                 `gorethink:"sender_id" json:"sender_id"`
	Opts      map[string]interface{} `gorethink:"opts" json:"opts"`
	CreatedAt time.Time              `gorethink:"created_at" json:"created_at"`
}

type ClientMessage struct {
	Type     int    `json:"type"`
	AckMsgId string `json:"ack_msg_id"`
}

func NewMessage(channelId string, typ int, payload interface{}, senderId string, opts map[string]interface{}) *Message {
	return &Message{
		ChannelId: channelId, Type: typ, Payload: payload,
		SenderId: senderId, Opts: opts, CreatedAt: time.Now(),
	}
}

func GetChannelByMembers(members []string) (*Channel, error) {
	h := md5.New()
	io.WriteString(h, strings.Join(members, ","))
	hash := hex.EncodeToString(h.Sum(nil))

	channel, err := FindChannelByHash(hash)
	if err != nil {
		return channel, err
	}

	if channel == nil {
		log.Printf("channel for members: %#v not exist, create a new", members)
		channel, err = CreateChannel(hash, members)
	}

	return channel, err
}

func FindChannelByHash(hash string) (*Channel, error) {
	query := r.Table("channels").Filter(r.Row.Field("hash").Eq(hash)).Limit(1)
	res, err := query.Run(rdb)

	if err != nil || res.IsNil() {
		return nil, err
	}

	channel := &Channel{}
	res.One(&channel)

	return channel, nil
}

func CreateChannel(hash string, members []string) (*Channel, error) {
	if len(hash) <= 0 || len(members) <= 0 {
		return nil, errors.New("hash or members is empty")
	}

	channel := &Channel{Hash: hash, Members: members,
		Type: CHANNEL_TYPE_NORMAL, CreatedAt: time.Now()}

	res, err := r.Table("channels").Insert(channel).RunWrite(rdb)
	if err != nil {
		return nil, err
	}

	channel.Id = res.GeneratedKeys[0]

	conn := pool.Get()
	defer conn.Close()

	key := "cm:" + channel.Id
	_, err = conn.Do("sadd", redis.Args{}.Add(key).AddFlat(members)...)
	if err != nil {
		log.Printf("sadd(%s) err: %s", key, err)
	}

	log.Println("redis sadd", key, members)
	return channel, nil
}

func FindMessage(id string) (*Message, error) {
	res, err := r.Table("messages").Get(id).Run(rdb)
	if err != nil || res.IsNil() {
		return nil, err
	}

	message := &Message{}
	res.One(message)

	// hacking the sent_at, otherwise it will be as float64 and json encode as sth. like 1.429238904e+09
	payload := message.Payload.(map[string]interface{})

	switch payload["sent_at"].(type) {
	case float64:
		payload["sent_at"] = int64(payload["sent_at"].(float64))
		message.Payload = payload
	default:
	}

	return message, err
}

func GetMessagesByChannel(channelId string) ([]*Message, error) {
	query := r.Table("messages").Filter(r.Row.Field("channel_id").Eq(channelId))
	res, err := query.Run(rdb)
	//	log.Printf("res=%#v, err=%s", res, err)

	if err != nil || res.IsNil() {
		return nil, err
	}

	messages := []*Message{}
	res.All(&messages)

	return messages, nil
}

func CreateMessage(m *Message) (*Message, error) {
	m.CreatedAt = time.Now()
	res, err := r.Table("messages").Insert(m).RunWrite(rdb)
	if err != nil {
		return nil, err
	}

	m.Id = res.GeneratedKeys[0]
	return m, nil
}

func GetUserQueuedMessageIds(userId string) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()

	ls, err := redis.Strings(conn.Do("lrange", fmt.Sprintf("mq:%s", userId), 0, -1))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var ids []string
	for _, v := range ls {
		ids = append(ids, v)
	}

	return ids, nil
}

func GetUserQueuedMessages(userId string) ([]string, []*Message, error) {
	ids, err := GetUserQueuedMessageIds(userId)
	if err != nil {
		return nil, nil, err
	}

	messages := []*Message{}

	for i := 0; i < len(ids); i++ {
		message, _ := FindMessage(ids[i])
		// if err2 != nil {
		// 	continue
		// }
		messages = append(messages, message)
	}

	return ids, messages, nil
}
