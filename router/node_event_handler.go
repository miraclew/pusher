package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"github.com/bitly/go-nsq"
)

type NodeEventHandler struct {
	app *App
}

func (n *NodeEventHandler) HandleMessage(message *nsq.Message) error {
	log.Debug("HandleNodeEvent %s", string(message.Body))
	evt := &push.NodeEvent{}
	err := json.Unmarshal(message.Body, evt)
	if err != nil {
		log.Error("Bad NodeEvent: body=%s err=%s", string(message.Body), err.Error())
	}
	if evt.Event == push.NODE_EVENT_ONLINE {
		body := &push.NodeEventOnline{}
		err = json.Unmarshal(evt.Body, body)
		if err != nil {
			log.Error("Bad NodeEventOnline: body=%s err=%s", string(evt.Body), err.Error())
			return err
		}

		if body.IsOnline {
			log.Debug("NodeEventOnline: userId: %d online", body.UserId)
			err = app.router.processQueue(body.UserId)
			if err != nil {
				log.Error("processQueue error: %s", err.Error())
			}
		}
	}
	return nil
}

func (n *NodeEventHandler) LogFailedMessage(m *nsq.Message) {
	n.app.db.Exec(`INSERT INTO messages
		(nsqd_address, topic, channel, body, attempts, timestamp) VALUES
		(?, ?, ?, ?, ?, ?)`,
		m.NSQDAddress, TOPIC_NODE_EVENT, CHANNEL_ROUTER, string(m.Body), m.Attempts, m.Timestamp)

}
