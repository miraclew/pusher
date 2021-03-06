package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"github.com/gorilla/websocket"
)

var connections = make(map[int64]*websocket.Conn)

func AddConnection(userId int64, conn *websocket.Conn) {
	old, ok := connections[userId]
	if ok && old != conn {
		log.Debug("AddConnection close old connection: %d/%s", userId, old.RemoteAddr().String())
		old.Close()
	}

	connections[userId] = conn
	go OnClientOnline(userId, true)
}

func RemoveConnection(userId int64, conn *websocket.Conn) {
	toRemove, ok := connections[userId]
	if ok && toRemove == conn {
		delete(connections, userId)
		push.RemoveClient(userId)
	} else {
		log.Debug("RemoveConnection %d/%s but new connection exist", userId, conn.RemoteAddr().String())
	}

	// go OnClientOnline(userId, false)
}

func GetConnection(userId int64) *websocket.Conn {
	return connections[userId]
}

func OnClientOnline(userId int64, online bool) {
	var err error
	evt := push.NodeEvent{}
	evt.Event = push.NODE_EVENT_ONLINE
	evt.NodeId = app.options.nodeId
	body := push.NodeEventOnline{}
	body.IsOnline = online
	body.UserId = userId
	evt.Body, err = json.Marshal(body)
	if err != nil {
		log.Error("OnClientOnline error1: %d", err.Error())
		return
	}

	b, err := json.Marshal(evt)
	if err != nil {
		log.Error("OnClientOnline error2: %d", err.Error())
		return
	}

	for _, producer := range app.producers {
		err := producer.Publish("node-event", b)
		if err != nil {
			log.Error("OnClientOnline Publish error: %s", err.Error())
		}
	}
}
