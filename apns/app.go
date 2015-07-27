package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"fmt"
	"github.com/anachronistic/apns"
	"github.com/bitly/go-nsq"
	"os"
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
	prodCert         string
	prodKey          string
	sandboxCert      string
	sandboxKey       string
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
	cmd := &push.ApnsCmd{}
	err := json.Unmarshal(message.Body, cmd)
	if err != nil {
		log.Error("body malformed: body=%s err=%s", string(message.Body), err.Error())
	}

	return a.pushToDevice(cmd)
}

func (a *App) pushToDevice(cmd *push.ApnsCmd) error {
	log.Info("apns pushToDevice %#v", cmd)
	payload := apns.NewPayload()
	payload.Alert = cmd.Alert
	payload.Sound = "ping.aiff"
	payload.Badge = cmd.Length
	pn := apns.NewPushNotification()
	pn.DeviceToken = cmd.DeviceToken
	pn.AddPayload(payload)

	cert := a.options.prodCert
	key := a.options.prodKey
	gatewayUrl := "gateway.push.apple.com:2195"
	if a.options.sandbox {
		cert := a.options.sandboxCert
		key := a.options.sandboxKey
		gatewayUrl = "gateway.sandbox.push.apple.com:2195"
	}

	client := apns.NewClient(gatewayUrl, cert, key)
	resp := client.Send(pn)
	if !resp.Success {
		log.Info("apns msgId:%s err: %s", cmd.MsgId, resp.Error)
	} else {
		log.Info("apns msgId:%s success", cmd.MsgId)
	}
	return nil
}

func (a *App) Exit() {
	close(a.exitChan)
	a.waitGroup.Wait()
}
