package main

import (
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

type Loader struct {
	numOfClients  int
	numOfMessages int
	redisAddr     string
	apiBaseUrl    string
	nsqdTCPAddr   string
	wsUrl         string
	clients       []*Client
	pool          *redis.Pool
	producer      *nsq.Producer
}

func NewLoader(numOfClients int, numOfMessages int, redisAddr string, apiBaseUrl string, nsqdTCPAddr string, wsUrl string) *Loader {
	return &Loader{
		numOfClients:  numOfClients,
		numOfMessages: numOfMessages,
		redisAddr:     redisAddr,
		apiBaseUrl:    apiBaseUrl,
		nsqdTCPAddr:   nsqdTCPAddr,
		wsUrl:         wsUrl,
	}
}

func (l *Loader) Start() {
	l.setup()
	l.loadClients()
}

func (l *Loader) Stop() {
	log.Printf("Stoping loader...")
	for i := 0; i < len(l.clients); i++ {
		client := l.clients[i]
		client.Stop()
	}
	l.tearDown()
}

func (l *Loader) setup() {
	l.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", l.redisAddr)
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

	l.createProducers()
}

func (l *Loader) tearDown() {
}

func (l *Loader) loadClients() error {
	conn := l.pool.Get()
	defer conn.Close()

	startUserId := 900000
	for i := 0; i < l.numOfClients; i++ {
		userId := startUserId + i
		client := NewClient(userId, l)

		log.Printf("NewClient: %d", userId)
		// setup tokens
		var values = map[string]string{
			"user_id":     fmt.Sprintf("%d", userId),
			"device_type": "1",
			"version":     "2.4.0",
		}

		_, err := conn.Do("hmset", redis.Args{}.Add(fmt.Sprintf("token:%d", userId)).AddFlat(values))
		if err != nil {
			log.Println("set token error: ", err.Error())
			return err
		}

		l.clients = append(l.clients, client)
	}

	for i := 0; i < len(l.clients); i++ {
		client := l.clients[i]
		time.Sleep(1 * time.Millisecond)
		client.Start()
	}

	log.Println("Create new clients: ", len(l.clients))
	return nil
}

func (l *Loader) createProducers() error {
	cfg := nsq.NewConfig()
	var err error
	l.producer, err = nsq.NewProducer(l.nsqdTCPAddr, cfg)
	if err != nil {
		log.Fatalf("failed to create nsq.Producer - %s", err)
	}

	return err
}

func (l *Loader) publish(msg string) error {
	return l.producer.Publish("server", []byte(msg))
}
