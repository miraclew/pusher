package pusher

import (
	"errors"
	"fmt"
	"github.com/anachronistic/apns"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

const (
	TYPE_DIRECT  = 1
	TYPE_CHANNEL = 2
	TYPE_BULK    = 3
)

type Hub struct {
	connections map[int64]*websocket.Conn
}

var hub *Hub

func GetHub() *Hub {
	if hub == nil {
		hub = &Hub{
			make(map[int64]*websocket.Conn),
		}
	}

	return hub
}

func (h *Hub) AddConnection(userId int64, conn *websocket.Conn) {
	h.connections[userId] = conn
	h.processQueue(userId)
}

func (h *Hub) RemoveConnection(userId int64) {
	delete(h.connections, userId)
}

func (h *Hub) PushMsg(msg *Message) {
	log.Printf("pushMsg: %#v\n", msg)
	if msg.ChannelId == "0" {
		h.broadcast(msg, true)
		return
	}

	if msg.Type == TYPE_DIRECT {
		receiverId, err := strconv.ParseInt(msg.ChannelId, 10, 64)
		if err != nil {
			log.Printf("channelId is not valid: %s", err.Error())
		} else {
			h.toUsers(msg, []int64{receiverId})
		}
	} else if msg.Type == TYPE_CHANNEL {
		h.toChannel(msg, msg.ChannelId)
	} else {
		//receivers := []string(msg.Options["receivers"])
		//h.toUsers(msg, receivers)
		//FIXME
	}
}

func (h *Hub) broadcast(msg *Message, online bool) {
	for _, v := range h.connections {
		v.WriteJSON(msg.Payload)
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
		var length int
		length, err = h.processQueue(userId)
		log.Printf("ret length=%d, err=%#v", length, err)

		if v, ok := msg.Opts["apn_enable"]; ok && v.(bool) {
			if length > 0 || err != nil {
				go h.pushToIosDevice(userId, msg, length)
			}
		}
	}

	return nil
}

func (h *Hub) sendToUser(userId int64, msg *Message) (ok bool, err error) {
	conn, ok := h.connections[userId]
	if ok {
		err = conn.WriteJSON(msg.Payload)
		if err != nil {
			return false, err
		}
		return true, err
	} else {
		h.pushToQueue(userId, msg, false)
		return false, nil
	}
}

func (h *Hub) pushToQueue(userId int64, msg *Message, left bool) (length int, err error) {
	cmd := "rpush"
	if left {
		cmd = "lpush"
	}

	return redis.Int(pool.Get().Do(cmd, fmt.Sprintf("mq:%d", userId), msg.Id))
}

func (h *Hub) processQueue(userId int64) (length int, err error) {
	log.Println("processQueue: ", userId)

	length, err = redis.Int(pool.Get().Do("llen", fmt.Sprintf("mq:%d", userId)))

	if err != nil {
		return length, err
	}

	if length <= 0 {
		return 0, nil
	}

	var msgId string
	msgId, err = redis.String(pool.Get().Do("rpop", fmt.Sprintf("mq:%d", userId)))
	if err != nil {
		log.Println(err)
		return length, err
	}

	length -= 1
	var msg *Message
	msg, err = FindMessage(msgId)

	if err != nil {
		return length, err
	}

	if msg == nil {
		return length, errors.New(fmt.Sprintf("msgId: %s not found", msgId))
	}

	var ok bool
	ok, err = h.sendToUser(userId, msg)
	if !ok || err != nil {
		return length + 1, err
	}

	if length > 0 {
		log.Println("mq len: ", length)
		return h.processQueue(userId)
	}

	return 0, nil
}

func (h *Hub) pushToIosDevice(userId int64, msg *Message, length int) error {
	deviceToken, err := redis.String(pool.Get().Do("get", fmt.Sprintf("apn_u2t:%d", userId)))
	if err != nil {
		log.Println("pushToIosDevice error: ", err)
		return err
	}

	if deviceToken != "" {
		log.Printf("apns msgId:%s userId:%d len:%d deviceToken=%s", msg.Id, userId, length, deviceToken)
		payload := apns.NewPayload()
		payload.Alert = "你有一条新的消息"
		payload.Sound = "ping.aiff"
		payload.Badge = length
		if v, ok := msg.Opts["apn_alert"]; ok {
			payload.Alert = v
		}

		pn := apns.NewPushNotification()
		pn.DeviceToken = deviceToken
		pn.AddPayload(payload)

		envDir := "prod"
		gatewayUrl := "gateway.push.apple.com:2195"

		if apnsDev {
			envDir = "dev"
			gatewayUrl = "gateway.sandbox.push.apple.com:2195"
		}
		certificateFile := "cert/" + envDir + "/cert.pem"
		keyFile := "cert/" + envDir + "/key.unencrypted.pem"

		client := apns.NewClient(gatewayUrl, certificateFile, keyFile)
		resp := client.Send(pn)

		if !resp.Success {
			log.Printf("apns msgId:%s err: %s", msg.Id, resp.Error)
		} else {
			log.Printf("apns msgId:%s success", msg.Id)
		}
	} else {
		log.Printf("user: %d offline, and has no apns device token\n", userId)
	}

	return nil
}

func (h *Hub) toChannel(msg *Message, channelId string) error {
	log.Println("push to channel_id:", channelId)

	ls, err := redis.Strings(pool.Get().Do("smembers", fmt.Sprintf("cm:%s", channelId)))
	if err != nil {
		log.Println(err)
		return err
	}

	var users []int64
	for _, v := range ls {
		userId, _ := strconv.ParseInt(v, 10, 64)
		users = append(users, userId)
	}

	log.Println("toUsers: ", users)
	return h.toUsers(msg, users)
}
