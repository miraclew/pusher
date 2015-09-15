package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"net/http"
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

	client, err := push.AuthClient(token)
	if err != nil {
		log.Warning("Auth failed, protocol=%s token=%s, err: %s", conn.Subprotocol(), token, err.Error())
		return
	}

	userId := client.UserId
	log.Info("New client, node=%d v=%s/%s p=%s token=%s userId=%d/%s", app.options.nodeId, client.Version, client.DeviceTypeName(), conn.Subprotocol(), token, userId, conn.RemoteAddr().String())
	client.NodeId = app.options.nodeId
	err = client.Save()
	if err != nil {
		log.Error("Client save err: %s", err.Error())
		return
	}
	client.Touch(app.options.clientTimeout)

	AddConnection(userId, conn)
	// Reading loop, required
	for {
		_, b, err2 := conn.ReadMessage()
		if err2 != nil {
			log.Info("Disconnect %d/%s, %s", userId, conn.RemoteAddr().String(), err2.Error())
			RemoveConnection(userId, conn)
			break
		} else {
			data := string(b)
			if data == "p" || data == "ping" {
				client.Touch(app.options.clientTimeout)
				conn.WriteMessage(websocket.TextMessage, []byte("q"))
				log.Debug("Pong => %d/%s", userId, conn.RemoteAddr().String())
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

func handleAck(userId int64, msgId string) {
	log.Debug("HandleAck %d %s \n", userId, msgId)

	conn := app.redisPool.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do("lrem", fmt.Sprintf("mq:%d", userId), 0, msgId))
	if err != nil {
		log.Error("lrem error: %s \n", err.Error())
	}
}
