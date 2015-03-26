package pusher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

const (
	TYPE_PRIVATE = 1
	TYPE_GROUP   = 2
	TYPE_BULK    = 3
)

type Hub struct {
	connections map[int64]io.ReadWriteCloser
}

func (h *Hub) AddConnection(userId int64, conn io.ReadWriteCloser) {
	h.connections[userId] = conn
}

func (h *Hub) RemoveConnection(userId int64) {
	delete(h.connections, userId)
}

func (h *Hub) pushMsg(msg *Message) {
	log.Printf("pushMsg: %#v\n", msg)
	if msg.ChannelId == "0" {
		h.broadcast(msg, true)
		return
	}

	if msg.Type == TYPE_PRIVATE || msg.Type == TYPE_GROUP {
		h.toChannel(msg, msg.ChannelId)
	} else {
		//receivers := []string(msg.Options["receivers"])
		//h.toUsers(msg, receivers)
		//FIXME
	}
}

func (h *Hub) broadcast(msg *Message, online bool) {
	for _, v := range h.connections {
		payload, _ := json.Marshal(msg.Payload)
		log.Println(payload, v)
		v.Write(payload)
	}
}

func (h *Hub) toUsers(msg *Message, users []int64) error {
	skipSender := true
	for i := 0; i < len(users); i++ {
		userId := users[i]
		senderId, err := strconv.ParseInt(msg.SenderId, 10, 64)
		if err != nil {
			log.Println("Warn: msg senderId is not integer", msg)
			continue
		}
		if userId == senderId && skipSender {
			log.Println("skipSender: ", userId)
			continue
		}

		// push to queue
		h.pushToQueue(userId, msg, true)
	}

	return nil
}

func (h *Hub) sendToUser(userId int64, msg *Message) error {
	conn, ok := h.connections[userId]
	if ok {
		// send to user
		payload, _ := json.Marshal(msg.Payload)
		_, err := conn.Write(payload)
		if err != nil {
			return err
		}
	} else {
		h.pushToQueue(userId, msg, false)
	}

	return nil
}

func (h *Hub) pushToQueue(userId int64, msg *Message, left bool) (length int, err error) {
	cmd := "rpush"
	if left {
		cmd = "lpush"
	}

	res := redis.Cmd(cmd, fmt.Sprintf("mq:%d", userId), msg.Id)
	if res.Err != nil {
		log.Println(res.Err)
		return -1, res.Err
	}

	return res.Int()
}

func (h *Hub) processQueue(userId int64) (length int, err error) {
	res := redis.Cmd("rpop", fmt.Sprintf("mq:%d", userId))
	if res.Err != nil {
		log.Println(res.Err)
		return -1, res.Err
	}

	var msgId string
	msgId, err = res.Str()
	if err != nil {
		return -1, err
	}

	var msg *Message
	msg, err = FindMessage(msgId)

	if err != nil {
		return -1, err
	}

	if msg == nil {
		return -1, errors.New(fmt.Sprintf("msgId: %s not found", msgId))
	}

	err = h.sendToUser(userId, msg)
	if err != nil {
		return -1, err
	}

	res = redis.Cmd("llen", fmt.Sprintf("mq:%d", userId))
	if res.Err != nil {
		log.Println(res.Err)
		return -1, res.Err
	}

	return res.Int()
}

func (h *Hub) pushToIosDevice(userId int64, msg *Message, length int) error {
	res := redis.Cmd("get", fmt.Sprintf("apn_u2t:%d", userId))
	if res.Err != nil {
		log.Println(res.Err)
		return res.Err
	}

	if deviceToken, _ := res.Str(); deviceToken != "" {
		// TODO: push to iOS device
	} else {
		log.Printf("user: %d offline, and has no apns device token\n", userId)
	}

	return nil
}

func (h *Hub) toChannel(msg *Message, channelId string) error {
	log.Println("push to channel_id:", channelId)

	res := redis.Cmd("smembers", fmt.Sprintf("cm:%s", channelId))
	if res.Err != nil {
		log.Println(res.Err)
		return res.Err
	}

	ls, _ := res.List()

	var users []int64
	for _, v := range ls {
		userId, _ := strconv.ParseInt(v, 10, 64)
		users = append(users, userId)
	}

	return h.toUsers(msg, users)
}
