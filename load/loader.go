package main

import (
	"fmt"
	rds "github.com/fzzy/radix/redis"
	"log"
	"time"
)

var redis *rds.Client

type Loader struct {
	numOfClients  int
	numOfMessages int
	redisAddr     string
	serverAddr    string
	clients       []*Client
}

func NewLoader(numOfClients int, numOfMessages int, redisAddr string, serverAddr string) *Loader {
	return &Loader{
		numOfClients:  numOfClients,
		numOfMessages: numOfMessages,
		redisAddr:     redisAddr,
		serverAddr:    serverAddr,
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
	var err error
	redis, err = rds.DialTimeout("tcp", l.redisAddr, time.Duration(10)*time.Second)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("redis connected: %s", l.redisAddr)
}

func (l *Loader) tearDown() {
	redis.Close()
}

func (l *Loader) loadClients() error {
	startUserId := 900000
	for i := 0; i < l.numOfClients; i++ {
		userId := startUserId + i
		client := NewClient(userId, l)

		log.Printf("NewClient: %d", userId)
		// setup tokens
		res := redis.Cmd("hmset", fmt.Sprintf("token:%d", userId), "user_id", userId)
		if res.Err != nil {
			log.Println("set token error: ", res.Err)
			return res.Err
		}

		l.clients = append(l.clients, client)
	}

	for i := 0; i < len(l.clients); i++ {
		client := l.clients[i]
		time.Sleep(100 * time.Millisecond)
		client.Start()
	}

	log.Println("Create new clients: ", len(l.clients))
	return nil
}
