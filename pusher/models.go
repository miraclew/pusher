package pusher

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	r "github.com/dancannon/gorethink"
	"io"
	"strings"
	"time"
)

const (
	CHANNEL_TYPE_NORMAL = 1
	MESSAGE_TYPE_NORMAL = 1
)

type Channel struct {
	Id        string    `gorethink:"id,omitempty"`
	Hash      string    `gorethink:"hash"`
	Members   []string  `gorethink:"members"`
	Type      int       `gorethink:"type"`
	CreatedAt time.Time `gorethink:"created_at"`
}

type Message struct {
	Id        string      `gorethink:"id,omitempty"`
	ChannelId string      `gorethink:"channel_id"`
	Type      int         `gorethink:"type"`
	Payload   interface{} `gorethink:"payload"`
	SenderId  string      `gorethink:"sender_id"`
	Options   interface{} `gorethink:"options"`
	CreatedAt time.Time   `gorethink:"created_at"`
}

func NewMessage(channelId string, typ int, payload interface{}, senderId string, options interface{}) *Message {
	return &Message{
		ChannelId: channelId, Type: typ, Payload: payload,
		SenderId: senderId, Options: options, CreatedAt: time.Now(),
	}
}

func GetChannelByMembers(members []string) (*Channel, error) {
	h := md5.New()
	io.WriteString(h, strings.Join(members, ","))
	hash := hex.EncodeToString(h.Sum(nil))

	fmt.Println(hash)
	channel, err := FindChannelByHash(hash)
	if err != nil {
		return channel, err
	}

	if channel == nil {
		return CreateChannel(hash, members)
	}

	return channel, err
}

func FindChannelByHash(hash string) (*Channel, error) {
	query := r.Db("mercury").Table("channels").Filter(r.Row.Field("hash").Eq(hash)).Limit(1)
	res, err := query.Run(session)

	if err != nil || res.IsNil() {
		return nil, err
	}

	channel := &Channel{}
	res.One(&channel)

	return channel, nil
}

func CreateChannel(hash string, members []string) (*Channel, error) {
	channel := &Channel{Hash: hash, Members: members,
		Type: CHANNEL_TYPE_NORMAL, CreatedAt: time.Now()}

	res, err := r.Db("mercury").Table("channels").Insert(channel).RunWrite(session)
	if err != nil {
		return nil, err
	}

	channel.Id = res.GeneratedKeys[0]
	return channel, nil
}

func FindMessage(id string) (*Message, error) {
	res, err := r.Db("mercury").Table("messages").Get(id).Run(session)
	if err != nil || res.IsNil() {
		return nil, err
	}

	message := &Message{}
	res.One(message)
	return message, err
}

func CreateMessage(m *Message) (*Message, error) {
	res, err := r.Db("mercury").Table("messages").Insert(m).RunWrite(session)
	if err != nil {
		return nil, err
	}

	m.Id = res.GeneratedKeys[0]
	return m, nil
}
