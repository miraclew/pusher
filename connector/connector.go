package main

import (
	"github.com/gorilla/websocket"
)

var connections = make(map[int64]*websocket.Conn)

func AddConnection(userId int64, conn *websocket.Conn) {
	connections[userId] = conn
}

func RemoveConnection(userId int64) {
	delete(connections, userId)
}

func GetConnection(userId int64) *websocket.Conn {
	return connections[userId]
}
