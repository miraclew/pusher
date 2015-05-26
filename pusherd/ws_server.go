package main

import (
	"coding.net/miraclew/pusher/pusher"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"strings"
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
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(req.URL.String())
	token := req.URL.Query().Get("token")
	version := req.URL.Query().Get("v")

	userId, err := pusher.GetUserIdByToken(token)
	if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
		//log.Fatalln("lost redis connection") // panic does not exit the process, use log.Fatal instead
		log.Println("Error: lost redis connection, reconnecting..")
	}

	if err != nil || userId <= 0 {
		log.Printf("Auth failed, protocol=%s token=%s, err: %s\n", conn.Subprotocol(), token, err.Error())
		conn.Close()
		return
	}

	//conn.WriteJSON(map[string]interface{}{"welcome": "hello, you are connected to push service"})
	log.Printf("New connection, v=%s protocol=%s token=%s userId=%d\n", version, conn.Subprotocol(), token, userId)

	pusher.GetHub().RemoveConnection(userId)
	pusher.GetHub().AddConnection(userId, conn)
	// Reading loop, required
	for {
		var msg = &pusher.ClientMessage{}
		err = conn.ReadJSON(msg)
		if err != nil {
			log.Println("Disconnect ", userId, err.Error())
			conn.Close()
			pusher.GetHub().RemoveConnection(userId)
			break
		} else {
			pusher.GetHub().HandleAck(userId, msg.AckMsgId)
		}
	}
}
