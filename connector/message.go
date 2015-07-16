package main

import (
	"time"
)

const (
	CHANNEL_TYPE_NORMAL = 1
	MESSAGE_TYPE_NORMAL = 1

	MSG_TYPE_ACK = 6001
)

type Message struct {
	Id        string      `json:"id"`
	ChannelId string      `json:"channel_id"`
	Type      int         `json:"type"`
	Payload   interface{} `json:"payload"`
	SenderId  string      `json:"sender_id"`
	// Opts      *MsgSendOpts `json:"opts"`
	CreatedAt time.Time `json:"created_at"`
}

type ClientMessage struct {
	Type     int    `json:"type"`
	AckMsgId string `json:"ack_msg_id"`
}
