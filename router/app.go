package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

type App struct {
	options   *AppOptions
	waitGroup sync.WaitGroup
	exitChan  chan int
	redisPool *redis.Pool
	db        *sqlx.DB
	consumer  *nsq.Consumer
	router    *Router
}

type AppOptions struct {
	redisAddr        string
	nsqdTCPAddrs     push.StringArray
	lookupdHTTPAddrs push.StringArray
	mysqlAddr        string
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

	db, err := sqlx.Connect("mysql", options.mysqlAddr)
	if err != nil {
		log.Error("connect mysql error: %s", err.Error())
		return nil
	}

	push.SetDb(db)
	push.SetRedisPool(pool)

	a := &App{
		options:   options,
		exitChan:  make(chan int),
		db:        db,
		redisPool: pool,
	}

	return a
}

func NewAppOptions() *AppOptions {
	options := &AppOptions{}

	return options
}

func (a *App) Main() {
	a.router = NewRouter()

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

	return a.router.route(&v)
}

func (a *App) Exit() {
	if a.consumer != nil {
		a.consumer.Stop()
	}
	close(a.exitChan)
	a.waitGroup.Wait()
}
