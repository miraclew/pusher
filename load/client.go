package main

type Client struct {
	UserId int64
}

func NewClient(userId int64) *Client {
	return &Client{userId}
}

func (c *Client) Start() {

}

func (c *Client) Stop() {

}
