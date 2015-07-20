package main

import (
	"coding.net/miraclew/pusher/app"
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
}

type AppOptions struct {
	redisAddr        string
	nsqdTCPAddrs     app.StringArray
	lookupdHTTPAddrs app.StringArray
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
	go a.startPubSub()
}

func (a *App) startPubSub() {
	cfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("topic", "channel", cfg)
	if err != nil {
		log.Error("nsq.NewConsumer error: %s", err.Error())
		panic(fmt.Sprintf("nsq.NewConsumer error: %s", err.Error()))
	}
	consumer.ConnectToNSQDs(a.options.nsqdTCPAddrs)
	consumer.ConnectToNSQLookupds(a.options.lookupdHTTPAddrs)
	consumer.AddHandler(a)
}

func (a *App) HandleMessage(message *nsq.Message) error {
	log.Debug("HandleMessage %#v", message)

	return nil
}

func (a *App) Exit() {

	close(a.exitChan)
	a.waitGroup.Wait()
}
