package main

import (
	// "coding.net/miraclew/pusher/pusher"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	urlString := fmt.Sprintf("ws://%s/ws?token=%d", c.Loader.serverAddr, c.UserId)

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
		log.Println(c.UserId, "writeLoop")
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
		log.Println(c.UserId, "readLoop")

		var v interface{}
		err := c.Conn.ReadJSON(&v)
		if err != nil {
			log.Printf("%d recv error:%s", c.UserId, err.Error())
			return
		}

		c.recvNum++
		log.Printf("%d recv:%#v", c.UserId, v)
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

	tpl := `{
		"type":1,
		"sender_id":"%d",
		"channel_id":"%d",
		"payload":{
			"id":"JmeUCZcjynsZyjpx",
			"type":4,
			"sub_type":4002,
			"chat_id":"0",
			"sender_id":"100054",
			"ttl":0,
			"sent_at":1429708702,
			"body":{"post":{"id":22,"text":"\u6211\u4eec\u7684\u751f\u6d3b","images":[],"audio":{"url":"","length":0}},"comment":"\u6d51\u8eab\u89e3\u6570"},
			"extra":{"sender_name":"The","sender_avatar":"http:\/static.uwang.me\/resource\/newavatar\/body_avatar_201.png","sender_vavatar":"http:\/\/static.uwang.me\/resource\/newavatar\/body201.png"}
		},
		"opts":{"ttl":0,"offlineEnable":true,"apnEnable":false,"alert":"\u4e00\u6761\u65b0\u6d88\u606f","apn_alert":"\u4f60\u6536\u5230\u4e00\u6761\u6d88\u606f"}
	}`
	var msg = fmt.Sprintf(tpl, c.UserId, recieverId)
	log.Println(msg)

	body := ioutil.NopCloser(strings.NewReader(msg))
	urlString := fmt.Sprintf("http://%s/direct_msg", c.Loader.serverAddr)
	resp, err := http.Post(urlString, "", body)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
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
