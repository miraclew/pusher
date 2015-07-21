package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/codegangsta/negroni"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/pat"
	"net/http"
	"sync"
	"time"
)

type App struct {
	options   *AppOptions
	waitGroup sync.WaitGroup
	exitChan  chan int
	redisPool *redis.Pool
	consumer  *nsq.Consumer
}

type AppOptions struct {
	wsIp             string
	wsPort           int
	nodeId           int
	redisAddr        string
	mysqlAddr        string
	nsqdTCPAddrs     push.StringArray
	lookupdHTTPAddrs push.StringArray
	apnsDev          bool
}

func NewApp(options *AppOptions) *App {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", options.redisAddr)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	a := &App{
		options:   options,
		exitChan:  make(chan int),
		redisPool: pool,
	}

	return a
}

func NewAppOptions() *AppOptions {
	options := &AppOptions{}

	return options
}

func (a *App) Main() {
	a.startWS()
	a.startConsumer()
}

func (a *App) startConsumer() {
	cfg := nsq.NewConfig()
	var err error
	a.consumer, err = nsq.NewConsumer("server", "router", cfg)
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
	// log.Debug("HandleMessage %#v", message)
	log.Debug("HandleMessage %s", string(message.Body))
	var v push.Message
	err := json.Unmarshal(message.Body, &v)
	if err != nil {
		log.Error("body malformed: body=%s err=%s", string(message.Body), err.Error())
		return err
	}

	return nil
}

func (a *App) startWS() {
	p := pat.New()
	p.Get("/ws", WSHandler)
	p.Get("/", WSHandler)

	n := negroni.Classic()
	n.UseHandler(p)

	go func() {
		addr := fmt.Sprintf("%s:%d", a.options.wsIp, a.options.wsPort+a.options.nodeId)
		log.Info("WebSocket listen: %s", addr)
		err := http.ListenAndServe(addr, n)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

// func (a *App) startPubSub() {
// 	conn := a.redisPool.Get()
// 	psc := redis.PubSubConn{conn}
// 	channel := fmt.Sprintf("nc:%d", a.options.nodeId)
// 	log.Info("redis subscribe %s", channel)
// 	psc.Subscribe(channel) // node channel
// 	for {
// 		switch v := psc.Receive().(type) {
// 		case redis.Message:
// 			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
// 		case redis.Subscription:
// 			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
// 		case error:
// 			log.Error("Pubsub receive error: %s", v.Error())
// 			time.Sleep(time.Millisecond)
// 		}
// 	}
// }

func (a *App) Exit() {

	close(a.exitChan)
	a.waitGroup.Wait()
}
