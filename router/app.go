package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"github.com/bitly/go-nsq"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
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
	a.nodeEventConsumer, err = nsq.NewConsumer("node-event", "router", cfg)
	if err != nil {
		log.Error("nsq.NewConsumer error: %s", err.Error())
		return err
	}
	a.nodeEventConsumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		log.Debug("HandleNodeEvent %s", string(m.Body))
		evt := &push.NodeEvent{}
		err := json.Unmarshal(m.Body, evt)
		if err != nil {
			log.Error("body malformed: body=%s err=%s", string(m.Body), err.Error())
		}
		if evt.Event == push.NODE_EVENT_ONLINE {
			body := &push.NodeEventOnline{}
			err = json.Unmarshal(evt.Body, body)
			if err != nil {
				log.Error("node evnt body malformed: body=%s err=%s", string(evt.Body), err.Error())
				return err
			}

			if body.IsOnline {
				log.Debug("NodeEventOnline: userId: %d online", body.UserId)
				err = a.router.processQueue(body.UserId)
				if err != nil {
					log.Error("processQueue error: %s", err.Error())
				}
			}
		}
		return nil
	}))

	a.nodeEventConsumer.ConnectToNSQDs(a.options.nsqdTCPAddrs)
	log.Info("ConnectToNSQDs %s", a.options.nsqdTCPAddrs.String())
	a.nodeEventConsumer.ConnectToNSQLookupds(a.options.lookupdHTTPAddrs)
	log.Info("ConnectToNSQLookupds %s", a.options.lookupdHTTPAddrs.String())
	return nil
}

func (a *App) startServerConsumer() error {
	cfg := nsq.NewConfig()
	var err error
	a.serverConsumer, err = nsq.NewConsumer("server", "router", cfg)
	if err != nil {
		log.Error("nsq.NewConsumer error: %s", err.Error())
		return err
	}

	a.serverConsumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		log.Debug("HandleServerMessage %s", string(m.Body))
		var v push.Message
		err := json.Unmarshal(m.Body, &v)
		if err != nil {
			log.Error("body malformed: body=%s err=%s", string(m.Body), err.Error())
			return err
		}

		return a.router.route(&v)
	}))

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
