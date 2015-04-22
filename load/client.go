package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	UserId int
	Loader *Loader
	Conn   *websocket.Conn
	quit   chan int
}

func NewClient(userId int, loader *Loader) *Client {
	client := &Client{}
	client.UserId = userId
	client.Loader = loader
	client.quit = make(chan int)
	return client
}

func (c *Client) Start() {
	conn, _, err := websocket.DefaultDialer.Dial(c.Loader.serverAddr, nil)
	if err != nil {
		log.Println(err)
	}

	c.Conn = conn
	go c.readLoop()
	go c.writeLoop()
}

func (c *Client) writeLoop() {
	for {
		c.Conn.WriteJSON(fmt.Sprintf("hello, this is %d", c.UserId))
		time.Sleep(1 * time.Second)
		// FIXME: waiting for a quit chan signal to exit
	}
}

func (c *Client) readLoop() {
	for {
		var v interface{}
		err := c.Conn.ReadJSON(&v)
		if err != nil {
			return
		}

		log.Printf("%d recv:%#v", c.UserId, v)
	}
}

func (c *Client) Stop() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
