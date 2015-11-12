package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/timehop/apns"
	"sync"
)

const (
	TOPIC_NAME   = "apns"
	CHANNEL_NAME = "worker"
)

type App struct {
	options   *AppOptions
	waitGroup sync.WaitGroup
	exitChan  chan int
	consumer  *nsq.Consumer
	producers map[string]*nsq.Producer
	client    *apns.Client
}

type AppOptions struct {
	nsqdTCPAddrs     push.StringArray
	lookupdHTTPAddrs push.StringArray
	sandbox          bool
	prodCert         string
	prodKey          string
	sandboxCert      string
	sandboxKey       string
}

func NewApp(options *AppOptions) *App {
	cert := options.prodCert
	key := options.prodKey
	gatewayUrl := "gateway.push.apple.com:2195"
	if options.sandbox {
		cert = options.sandboxCert
		key = options.sandboxKey
		gatewayUrl = "gateway.sandbox.push.apple.com:2195"
	}

	log.Info("create client: %s %s %s", gatewayUrl, cert, key)
	client, err := apns.NewClientWithFiles(gatewayUrl, cert, key)
	if err != nil {
		log.Fatalf("could not create new client: %s", err.Error())
	}

	go func() {
		for f := range client.FailedNotifs {
			log.Error("Notif", f.Notif.ID, "failed with", f.Err.Error())
		}
	}()

	a := &App{
		options:  options,
		exitChan: make(chan int),
		client:   &client,
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
	a.consumer, err = nsq.NewConsumer(TOPIC_NAME, CHANNEL_NAME, cfg)
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
	cmd := &push.ApnsCmd{}
	err := json.Unmarshal(message.Body, cmd)
	if err != nil {
		log.Error("body malformed: body=%s err=%s", string(message.Body), err.Error())
	}

	return a.pushToDevice(cmd)
}

func (a *App) LogFailedMessage(m *nsq.Message) {
	log.Critical("LogFailedMessage(%s, %s, %s, %s, %d, %d)",
		m.NSQDAddress, TOPIC_NAME, CHANNEL_NAME, string(m.Body), m.Attempts, m.Timestamp)
}

func (a *App) pushToDevice(cmd *push.ApnsCmd) error {
	log.Info("apns pushToDevice %#v", cmd)
	payload := apns.NewPayload()
	payload.APS.Alert.Title = cmd.Alert
	payload.APS.Sound = "ping.aiff"
	// payload.Badge = cmd.Length
	badge := 1
	payload.APS.Badge = &badge
	if cmd.Payload != nil {
		payload.APS.Alert.Body = string(cmd.Payload)
	}

	pn := apns.NewNotification()
	pn.DeviceToken = cmd.DeviceToken
	pn.Priority = apns.PriorityImmediate
	pn.Payload = payload
	pn.ID = cmd.MsgId

	log.Debug("apns.Client %#v send payload: %#v", a.client, pn)
	a.client.Send(pn)
	log.Info("apns send msgId:%s", cmd.MsgId)

	return nil
}

func (a *App) Exit() {
	close(a.exitChan)
	a.waitGroup.Wait()
}
