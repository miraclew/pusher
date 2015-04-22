package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Client struct {
	UserId int
	Loader *Loader
	Conn   *websocket.Conn
}

func NewClient(userId int, loader *Loader) *Client {
	return &Client{userId, loader}
}

func (c *Client) Start() {
	var err error
	var response *http.Response
	c.Conn, response, err = websocket.DefaultDialer.Dial(c.Loader.serverAddr, nil)
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) Stop() {

}
