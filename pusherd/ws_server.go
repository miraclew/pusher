package main

import (
	"coding.net/miraclew/pusher/pusher"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type Bus interface {
	AddConnection(int64, io.ReadWriteCloser)
	RemoveConnection(int64)
}

type WsServer struct {
	doneCh chan bool
	errCh  chan error
	bus    Bus
}

var wsServer *WsServer

func NewWsServer(bus Bus) *WsServer {
	if wsServer == nil {
		wsServer = &WsServer{
			make(chan bool),
			make(chan error),
			bus,
		}
	}

	return wsServer
}

func (w *WsServer) onConnected(conn *websocket.Conn) {
	var userId int64
	token := conn.Request().URL.Query().Get("token")
	userId, err := pusher.GetUserIdByToken(token)
	//conn.Write([]byte("Token invalid, Good bye!"))
	if err != nil || userId <= 0 {
		conn.Close()
	}

	w.bus.AddConnection(userId, conn)
	log.Printf("New connection, %s -> %d \n", token, userId)
}

func WsServe(listener net.Listener, bus Bus) {
	log.Printf("WS: listening on %s", listener.Addr().String())
	s := NewWsServer(bus)

	wsHandler := &websocket.Server{Handler: websocket.Handler(s.onConnected)}

	httpServer := &http.Server{
		Handler: wsHandler,
	}

	err := httpServer.Serve(listener)
	// theres no direct way to detect this error because it is not exposed
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		log.Printf("ERROR: ws.Serve() - %s", err.Error())
	}

	log.Printf("HTTP: closing %s", listener.Addr().String())
}
