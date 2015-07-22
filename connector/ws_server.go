package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var clients map[int64]*push.Client

func init() {
	clients = make(map[int64]*push.Client)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
	Subprotocols:    []string{"gx-v1", "gx-v2"},
}

func WSHandler(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil || conn == nil {
		log.Error(err.Error())
		return
	}

	defer conn.Close()

	token := req.URL.Query().Get("token")
	if token == "" || token == "null" {
		return
	}

	client, err := authClient(token)

	if err != nil {
		log.Warning("Auth failed, protocol=%s token=%s, err: %s", conn.Subprotocol(), token, err.Error())
		return
	}

	userId := client.UserId
	log.Info("New client, v=%s/%s p=%s token=%s userId=%d", client.Version, client.DeviceTypeName(), conn.Subprotocol(), token, userId)

	RemoveConnection(userId)
	AddConnection(userId, conn)
	// Reading loop, required
	for {
		_, b, err2 := conn.ReadMessage()
		if err2 != nil {
			log.Info("Disconnect ", userId, err2.Error())
			RemoveConnection(userId)
			break
		} else {
			data := string(b)
			if data == "p" || data == "ping" {
				conn.WriteMessage(websocket.TextMessage, []byte("q"))
				continue
			}

			var msg = &push.ClientMessage{}
			err = json.Unmarshal(b, msg)
			if err != nil {
				log.Error("Malformed msg: %d %s %s", userId, data, err.Error())
			}

			if msg.Type == push.MSG_TYPE_ACK {
				handleAck(userId, msg.AckMsgId)
			}
		}
	}
}

func authClient(token string) (*push.Client, error) {
	conn := app.redisPool.Get()
	defer conn.Close()

	v, err := redis.StringMap(conn.Do("hgetall", "token:"+token))

	if err != nil {
		return nil, err
	}
	log.Debug("token:%s %#v", token, v)

	userId, err := strconv.ParseInt(v["user_id"], 10, 64)
	if err != nil {
		return nil, err
	}

	if v["device_type"] == "" {
		v["device_type"] = "2"
	}

	deviceType, err := strconv.ParseInt(v["device_type"], 10, 64)
	if err != nil {
		return nil, err
	}

	client := &push.Client{}
	// client.Token = token
	client.UserId = userId
	client.Version = v["version"]
	client.DeviceType = int(deviceType)
	client.NodeId = app.options.nodeId

	clients[client.UserId] = client
	client.Save()

	return client, nil
}

func handleAck(userId int64, msgId string) {
	log.Debug("HandleAck %d %s \n", userId, msgId)

	conn := app.redisPool.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("lrem", fmt.Sprintf("mq:%d", userId), 0, msgId))
	if err != nil {
		log.Error("lrem error: %s \n", err.Error())
	}
}
