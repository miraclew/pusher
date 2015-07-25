package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"fmt"
	"github.com/anachronistic/apns"
	"github.com/bitly/go-nsq"
	"sync"
)

type App struct {
	options   *AppOptions
	waitGroup sync.WaitGroup
	exitChan  chan int
	consumer  *nsq.Consumer
	producers map[string]*nsq.Producer
}

type AppOptions struct {
	nsqdTCPAddrs     push.StringArray
	lookupdHTTPAddrs push.StringArray
	sandbox          bool
}

func NewApp(options *AppOptions) *App {
	a := &App{
		options:  options,
		exitChan: make(chan int),
	}

	return a
}

func NewAppOptions() *AppOptions {
	options := &AppOptions{}

	return options
}

func (a *App) Main() {
	a.createProducers()
	a.startConsumer()
}

func (a *App) createProducers() {
	a.producers = make(map[string]*nsq.Producer)
	cfg := nsq.NewConfig()
	// cfg.UserAgent = fmt.Sprintf("to_nsq/%s go-nsq/%s", version.Binary, nsq.VERSION)

	for _, addr := range a.options.nsqdTCPAddrs {
		producer, err := nsq.NewProducer(addr, cfg)
		if err != nil {
			log.Fatalf("failed to create nsq.Producer - %s", err)
		}
		a.producers[addr] = producer
	}

	if len(a.producers) == 0 {
		log.Fatal("--nsqd-tcp-address required")
	}
}

func (a *App) startConsumer() {
	cfg := nsq.NewConfig()
	var err error
	a.consumer, err = nsq.NewConsumer("apns", "worker", cfg)
	if err != nil {
		log.Error("nsq.NewConsumer error: %s", err.Error())
		panic(fmt.Sprintf("nsq.NewConsumer error: %s", err.Error()))
	}
	a.consumer.AddHandler(a)

	a.consumer.ConnectToNSQDs(a.options.nsqdTCPAddrs)
	log.Info("ConnectToNSQDs %s", a.options.nsqdTCPAddrs.String())
	a.consumer.ConnectToNSQLookupds(a.options.lookupdHTTPAddrs)
	log.Info("ConnectToNSQLookupds %s", a.options.lookupdHTTPAddrs.String())
}

func (a *App) HandleMessage(message *nsq.Message) error {
	// log.Debug("HandleMessage %s", string(message.Body))
	cmd := &push.ApnsCmd{}
	err := json.Unmarshal(message.Body, cmd)
	if err != nil {
		log.Error("body malformed: body=%s err=%s", string(message.Body), err.Error())
	}

	// log.Debug("NodeCmdPush: msgId=%d receiverId: %d payload: %s", body.MsgId, body.ReceiverId, string(body.Payload))

	return nil
}

func (a *App) pushToIosDevice(userId int64, msgId string, deviceToken string, alert string, length int) error {
	log.Info("apns msgId:%s userId:%d len:%d deviceToken=%s", msgId, userId, length, deviceToken)
	payload := apns.NewPayload()
	payload.Alert = alert
	payload.Sound = "ping.aiff"
	payload.Badge = length
	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	envDir := "prod"
	gatewayUrl := "gateway.push.apple.com:2195"
	if a.options.sandbox {
		envDir = "dev"
		gatewayUrl = "gateway.sandbox.push.apple.com:2195"
	}
	certificateFile := "cert/" + envDir + "/cert.pem"
	keyFile := "cert/" + envDir + "/key.unencrypted.pem"
	client := apns.NewClient(gatewayUrl, certificateFile, keyFile)
	resp := client.Send(pn)
	if !resp.Success {
		log.Info("apns msgId:%s err: %s", msgId, resp.Error)
	} else {
		log.Info("apns msgId:%s success", msgId)
	}
	return nil
}

func (a *App) Exit() {

	close(a.exitChan)
	a.waitGroup.Wait()
}
