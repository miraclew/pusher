package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
)

type Client struct {
	token string
}

func _main() {
	fmt.Println("start")
	origin := "http://localhost/"
	url := "ws://localhost:9010/ws"
	ws, err := websocket.Dial(url, "gx-v1", origin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}
