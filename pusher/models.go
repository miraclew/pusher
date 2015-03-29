package pusher

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	r "github.com/dancannon/gorethink"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	CHANNEL_TYPE_NORMAL = 1
	MESSAGE_TYPE_NORMAL = 1
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
	Options   map[string]interface{} `gorethink:"options" json:"options"`
	CreatedAt time.Time              `gorethink:"created_at" json:"created_at"`
}

func NewMessage(channelId string, typ int, payload interface{}, senderId string, options map[string]interface{}) *Message {
	return &Message{
		ChannelId: channelId, Type: typ, Payload: payload,
		SenderId: senderId, Options: options, CreatedAt: time.Now(),
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
		channel, err = CreateChannel(hash, members)
	}

	return channel, err
}

func FindChannelByHash(hash string) (*Channel, error) {
	query := r.Db("mercury").Table("channels").Filter(r.Row.Field("hash").Eq(hash)).Limit(1)
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

	res, err := r.Db("mercury").Table("channels").Insert(channel).RunWrite(rdb)
	if err != nil {
		return nil, err
	}

	channel.Id = res.GeneratedKeys[0]

	key := "cm:" + channel.Id
	res2 := redis.Cmd("sadd", key, members)
	if res2.Err != nil {
		log.Printf("sadd(%s) err: %s", key, res2.Err)
	}

	log.Println("redis sadd", key, members)
	return channel, nil
}

func FindMessage(id string) (*Message, error) {
	res, err := r.Db("mercury").Table("messages").Get(id).Run(rdb)
	if err != nil || res.IsNil() {
		return nil, err
	}

	message := &Message{}
	res.One(message)
	return message, err
}

func CreateMessage(m *Message) (*Message, error) {
	res, err := r.Db("mercury").Table("messages").Insert(m).RunWrite(rdb)
	if err != nil {
		return nil, err
	}

	m.Id = res.GeneratedKeys[0]
	return m, nil
}

func GetUserIdByToken(token string) (int64, error) {
	res := redis.Cmd("hgetall", "token:"+token)

	if res.Err != nil {
		return 0, res.Err
	}

	h, err := res.Hash()
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(h["user_id"], 10, 64)
}
