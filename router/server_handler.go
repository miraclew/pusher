package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"github.com/bitly/go-nsq"
)

type ServerHandler struct {
	app *App
}

func (s *ServerHandler) HandleMessage(message *nsq.Message) error {
	log.Debug("HandleServerMessage %s", string(message.Body))
	var v push.Message
	err := json.Unmarshal(message.Body, &v)
	if err != nil {
		log.Error("Bad Push.Message: body=%s err=%s", string(message.Body), err.Error())
		return nil // JUST FIN the message if it's not good formed
	}

	return app.router.route(&v)
}

func (s *ServerHandler) LogFailedMessage(m *nsq.Message) {
	s.app.db.Exec(`INSERT INTO messages
		(nsqd_address, topic, channel, body, attempts, timestamp) VALUES
		(?, ?, ?, ?, ?, ?)`,
		m.NSQDAddress, TOPIC_SERVER, CHANNEL_ROUTER, string(m.Body), m.Attempts, m.Timestamp)
}
