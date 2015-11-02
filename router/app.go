package main

import (
	"coding.net/miraclew/pusher/push"
	"github.com/bitly/go-nsq"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

const (
	TOPIC_SERVER     = "server"
	CHANNEL_ROUTER   = "router"
	TOPIC_NODE_EVENT = "node-event"
)

type App struct {
	options           *AppOptions
	waitGroup         sync.WaitGroup
	exitChan          chan int
	redisPool         *redis.Pool
	db                *sqlx.DB
	serverConsumer    *nsq.Consumer
	nodeEventConsumer *nsq.Consumer
	router            *Router
}

type AppOptions struct {
	redisAddr        string
	nsqdTCPAddrs     push.StringArray
	lookupdHTTPAddrs push.StringArray
	nsqTopic         string
	nsqChannel       string
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
	a.startNodeEventConsumer()
	a.startServerConsumer()
}

func (a *App) startNodeEventConsumer() error {
	cfg := nsq.NewConfig()
	var err error
	a.nodeEventConsumer, err = nsq.NewConsumer(TOPIC_NODE_EVENT, CHANNEL_ROUTER, cfg)
	if err != nil {
		log.Error("nsq.NewConsumer error: %s", err.Error())
		return err
	}

	handler := &NodeEventHandler{app: a}
	a.nodeEventConsumer.AddHandler(handler)

	a.nodeEventConsumer.ConnectToNSQDs(a.options.nsqdTCPAddrs)
	log.Info("ConnectToNSQDs %s", a.options.nsqdTCPAddrs.String())
	a.nodeEventConsumer.ConnectToNSQLookupds(a.options.lookupdHTTPAddrs)
	log.Info("ConnectToNSQLookupds %s", a.options.lookupdHTTPAddrs.String())
	return nil
}

func (a *App) startServerConsumer() error {
	cfg := nsq.NewConfig()
	cfg.MaxBackoffDuration = time.Second

	var err error
	a.serverConsumer, err = nsq.NewConsumer(a.options.nsqTopic, CHANNEL_ROUTER, cfg)
	if err != nil {
		log.Error("nsq.NewConsumer error: %s", err.Error())
		return err
	}

	handler := &ServerHandler{app: a}
	a.serverConsumer.AddHandler(handler)

	a.serverConsumer.ConnectToNSQDs(a.options.nsqdTCPAddrs)
	log.Info("ConnectToNSQDs %s", a.options.nsqdTCPAddrs.String())
	a.serverConsumer.ConnectToNSQLookupds(a.options.lookupdHTTPAddrs)
	log.Info("ConnectToNSQLookupds %s", a.options.lookupdHTTPAddrs.String())
	return nil
}

func (a *App) Exit() {
	if a.serverConsumer != nil {
		a.serverConsumer.Stop()
	}
	close(a.exitChan)
	a.waitGroup.Wait()
}
