package main

import (
	"coding.net/miraclew/pusher/pusher"
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
	if err != nil {
		log.Println(err)
		return
	}

	token := req.URL.Query().Get("token")
	conn.WriteJSON(map[string]interface{}{"welcome": "bob"})

	userId, err := pusher.GetUserIdByToken(token)
	if err != nil || userId <= 0 {
		log.Printf("Auth failed, protocol=%s token=%s\n", conn.Subprotocol(), token)
		conn.Close()
	}

	pusher.GetHub().AddConnection(userId, conn)
	log.Printf("New connection, protocol=%s token=%s\n", conn.Subprotocol(), token)

	// Reading loop, required
	for {
		if _, _, err := conn.NextReader(); err != nil {
			conn.Close()
			break
		}
	}
}
