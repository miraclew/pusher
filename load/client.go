package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	UserId      int
	Loader      *Loader
	Conn        *websocket.Conn
	wQuit       chan int
	stopped     bool
	sendNum     int
	recvNum     int
	recieverIdx int
}

func NewClient(userId int, loader *Loader) *Client {
	client := &Client{}
	client.UserId = userId
	client.Loader = loader
	client.sendNum = 0
	client.recvNum = 0
	client.stopped = true
	return client
}

func (c *Client) Start() {
	urlString := fmt.Sprintf("%s?token=%d", c.Loader.wsUrl, c.UserId)

	conn, _, err := websocket.DefaultDialer.Dial(urlString, nil)
	if err != nil {
		log.Printf("Dail %s error: %s", urlString, err)
		return
	}

	log.Printf("%d connected", c.UserId)
	c.stopped = false

	c.Conn = conn
	go c.readLoop()
	go c.writeLoop()
}

func (c *Client) writeLoop() {
	c.wQuit = make(chan int)

	for {
		timer := time.After(2 * time.Second)
		// log.Println(c.UserId, "writeLoop")
		select {
		case <-c.wQuit:
			log.Println(c.UserId, "writeLoop End")
			return
		case <-timer:
			c.sendMessage()
			c.sendNum++
		}
	}
}

func (c *Client) readLoop() {
	defer c.Stop()

	for {
		// log.Println(c.UserId, "readLoop")

		var v interface{}
		err := c.Conn.ReadJSON(&v)
		if err != nil {
			log.Printf("%d recv error:%s", c.UserId, err.Error())
			return
		}

		c.recvNum++
		// log.Printf("%d recv:%#v", c.UserId, v)
	}
}

func (c *Client) sendMessage() {
	if c.recieverIdx >= len(c.Loader.clients) {
		c.recieverIdx = 0
	}

	defer func() {
		c.recieverIdx++
	}()

	recieverId := c.Loader.clients[c.recieverIdx].UserId

	if recieverId == c.UserId {
		return
	}

	tpl := `{"type":1,
	"sub_type":1001,
	"chat_id":0,
	"sender_id":%d,
	"receiver":"%d",
	"body":"{\"mime\":\"audio\",\"content\":{\"url\":\"http:\\\/\\\/static2.uwang.me\\\/audio\\\/2015\\\/07\\\/27\\\/155b5a03734053.mp3\",\"length\":4}}",
	"opts":"{\"ttl\":0,\"offline_enable\":true,\"ack_enable\":true,\"apn_enable\":true,\"alert\":\"\"}",
	"extra":"{\"sender_name\":\"jvcol\",\"sender_avatar\":\"http:\\\/static.uwang.me\\\/resource\\\/newavatar\\\/body_avatar_205.png\",\"sender_vavatar\":\"http:\\\/\\\/static.uwang.me\\\/resource\\\/newavatar\\\/body205.png\",\"age\":25,\"love_status\":1,\"gender\":1,\"address\":\"\\u4e0a\\u6d77\\u6d66\\u4e1c\\u65b0\"}",
	"timestamp":1437966393768,
	"id":415
	}`

	var msg = fmt.Sprintf(tpl, c.UserId, recieverId)
	c.Loader.publish(msg)
}

func (c *Client) Stop() {
	if c.stopped {
		return
	}
	c.stopped = true

	if c.Conn != nil {
		c.Conn.Close()
	}

	if c.wQuit != nil {
		c.wQuit <- 0
	}
}
