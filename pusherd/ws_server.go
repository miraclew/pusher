package main

import (
	"coding.net/miraclew/pusher/pusher"
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

type ConnectionManager interface {
	AddConnection(int64, io.ReadWriteCloser)
	RemoveConnection(int64)
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
		log.Println(err)
		return
	}

	defer conn.Close()

	token := req.URL.Query().Get("token")
	if token == "" || token == "null" {
		return
	}

	client, err := pusher.AuthClient(token)

	if err != nil {
		log.Printf("Auth failed, protocol=%s token=%s, err: %s\n", conn.Subprotocol(), token, err.Error())
		return
	}

	userId := client.UserId
	log.Printf("New client, v=%s/%s p=%s token=%s userId=%d\n", client.Version, client.DeviceTypeName(), conn.Subprotocol(), token, userId)

	pusher.GetHub().RemoveConnection(userId)
	pusher.GetHub().AddConnection(userId, conn)
	// Reading loop, required
	for {
		_, b, err2 := conn.ReadMessage()
		if err2 != nil {
			log.Println("Disconnect ", userId, err2.Error())
			pusher.GetHub().RemoveConnection(userId)
			break
		} else {
			data := string(b)
			if data == "p" || data == "ping" {
				conn.WriteMessage(websocket.TextMessage, []byte("q"))
				continue
			}

			var msg = &pusher.ClientMessage{}
			err = json.Unmarshal(b, msg)
			if err != nil {
				log.Printf("Malformed msg: %d %s %s", userId, data, err.Error())
			}

			if msg.Type == pusher.MSG_TYPE_ACK {
				pusher.GetHub().HandleAck(userId, msg.AckMsgId)
			}
		}
	}
}
