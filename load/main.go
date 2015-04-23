package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	numOfClients  = flag.Int("n", 2, "Number of clients to start")
	numOfMessages = flag.Int("m", 5, "Number of messages to send per second")
	redisAddr     = flag.String("redisAddr", "127.0.0.1:6379", "<addr>:<port> (127.0.0.1:6379) redis address to connect")
	serverAddr    = flag.String("serverAddr", "localhost:9001", "API address")
)

func main() {
	flag.Parse()

	if *redisAddr == "" {
		fmt.Println("redisAddr is required")
		return
	}

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	loader := NewLoader(*numOfClients, *numOfMessages, *redisAddr, *serverAddr)
	loader.Start()

	<-exitChan
	loader.Stop()
}
