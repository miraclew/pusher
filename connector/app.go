package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/codegangsta/negroni"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
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
	producers map[string]*nsq.Producer
}

type AppOptions struct {
	wsIp             string
	wsPort           int
	wsPortStart      int
	nodeId           int
	redisAddr        string
	clientTimeout    int
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

	push.SetRedisPool(pool)

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
	a.createProducers()
	a.startConsumer()

	a.startWS()
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
	cfg.MaxBackoffDuration = time.Second
	cfg.MaxAttempts = 1
	cfg.MaxRequeueDelay = time.Second
	cfg.DefaultRequeueDelay = time.Second

	var err error
	a.consumer, err = nsq.NewConsumer(fmt.Sprintf("connector-%d", a.options.nodeId), "connector", cfg)
	if err != nil {
		log.Error("nsq.NewConsumer error: %s", err.Error())
		panic(fmt.Sprintf("nsq.NewConsumer error: %s", err.Error()))
	}
	a.consumer.AddConcurrentHandlers(a, 100)

	a.consumer.ConnectToNSQDs(a.options.nsqdTCPAddrs)
	log.Info("ConnectToNSQDs %s", a.options.nsqdTCPAddrs.String())
	a.consumer.ConnectToNSQLookupds(a.options.lookupdHTTPAddrs)
	log.Info("ConnectToNSQLookupds %s", a.options.lookupdHTTPAddrs.String())
}

func (a *App) HandleMessage(message *nsq.Message) error {
	// log.Debug("HandleMessage %s", string(message.Body))
	cmd := &push.NodeCmd{}
	err := json.Unmarshal(message.Body, cmd)
	if err != nil {
		log.Error("body malformed: body=%s err=%s", string(message.Body), err.Error())
		return nil
	}
	if cmd.Cmd == push.NODE_CMD_PUSH {
		body := &push.NodeCmdPush{}
		err = json.Unmarshal(cmd.Body, body)

		log.Debug("NodeCmdPush: msgId=%d receiverId: %d payload: %s", body.MsgId, body.ReceiverId, string(body.Payload))
		conn := GetConnection(body.ReceiverId)
		if conn != nil {
			err := conn.WriteMessage(websocket.TextMessage, body.Payload)
			if err != nil {
				log.Error("WriteMessage err: %s", err.Error())
				return nil
			} else {
				log.Info("Send OK msgId: %d => userId: %d", body.MsgId, body.ReceiverId)
			}
		} else {
			log.Warning("User %d is not connected", body.ReceiverId)
		}
	}

	return nil
}

func (a *App) LogFailedMessage(m *nsq.Message) {
	log.Critical("LogFailedMessage(%s, %s, %s, %s, %d, %d)",
		m.NSQDAddress, fmt.Sprintf("connector-%d", a.options.nodeId), "connector", string(m.Body), m.Attempts, m.Timestamp)
}

func (a *App) startWS() {
	p := pat.New()
	p.Get("/ws", WSHandler)
	p.Get("/", WSHandler)

	n := negroni.Classic()
	n.UseHandler(p)

	go func() {
		port := a.options.wsPort
		if port == 0 {
			port = a.options.wsPortStart + a.options.nodeId
		}
		addr := fmt.Sprintf("%s:%d", a.options.wsIp, port)
		log.Info("WebSocket listen: %s", addr)
		err := http.ListenAndServe(addr, n)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (a *App) Exit() {

	close(a.exitChan)
	a.waitGroup.Wait()
}
