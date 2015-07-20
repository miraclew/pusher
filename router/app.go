package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/garyburd/redigo/redis"
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
	redisAddr        string
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

func (a *App) Exit() {
	if a.consumer != nil {
		a.consumer.Stop()
	}
	close(a.exitChan)
	a.waitGroup.Wait()
}
