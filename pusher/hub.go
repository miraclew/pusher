package pusher

import (
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

func (h *Hub) ConnectionsCount() int {
	return len(h.connections)
}

func (h *Hub) PushMsg(msg *Message) {
	//log.Printf("pushMsg: %#v\n", msg)
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
		_, ok := h.connections[userId]
		if ok { // online
			err = h.processQueue(userId)
			if err != nil {
				log.Printf("processQueue error=%s", err.Error())
			}
		} else { // offline
			if v, ok := msg.Opts["apn_enable"]; ok && v.(bool) {
				go h.pushToIosDevice(userId, msg, 1) // TODO:
			}
		}
	}

	return nil
}

func (h *Hub) pushToQueue(userId int64, msg *Message, left bool) (length int, err error) {
	cmd := "rpush"
	if left {
		cmd = "lpush"
	}

	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do(cmd, fmt.Sprintf("mq:%d", userId), msg.Id))
}

func (h *Hub) processQueue(userId int64) (err error) {
	log.Println("processQueue: ", userId)

	ws, ok := h.connections[userId]
	if !ok {
		return nil
	}

	conn := pool.Get()
	defer conn.Close()
	ids, err2 := redis.Strings(conn.Do("lrange", fmt.Sprintf("mq:%d", userId), 0, 9))

	if err2 != nil {
		return err2
	}

	if len(ids) <= 0 {
		return nil
	}

	for i := 0; i < len(ids); i++ {
		msgId := ids[i]

		var msg *Message
		msg, err = FindMessage(msgId)

		if err != nil {
			log.Println(fmt.Sprintf("FindMessage: %s error: %s", msgId, err.Error()))
			continue
		}

		if msg == nil {
			log.Println(fmt.Sprintf("msgId: %s not found", msgId))
			continue
		}

		err = ws.WriteJSON(msg.Payload)
		if err != nil {
			log.Printf("Error: %d WriteJSON error: %s \n", userId, err.Error())
		}
	}

	return nil
}

func (h *Hub) pushToIosDevice(userId int64, msg *Message, length int) error {
	conn := pool.Get()
	defer conn.Close()

	deviceToken, err := redis.String(conn.Do("get", fmt.Sprintf("apn_u2t:%d", userId)))
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

	conn := pool.Get()
	defer conn.Close()

	ls, err := redis.Strings(conn.Do("smembers", fmt.Sprintf("cm:%s", channelId)))
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

func (h *Hub) HandleAck(userId int64, msgId string) {

}
