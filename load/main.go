package main

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

var (
	numOfClients  = flag.Int("n", 2, "Number of clients to start")
	numOfMessages = flag.Int("m", 5, "Number of messages to send per second")
)

func main() {
	flag.Parse()
	godotenv.Load()

	nsqdTCPAddr := os.Getenv("NSQD_ADDR")
	redisAddr := os.Getenv("REDIS_ADDR")
	apiBaseUrl := os.Getenv("API_BASE_URL")
	wsUrl := os.Getenv("WS_URL")

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	loader := NewLoader(*numOfClients, *numOfMessages, redisAddr, apiBaseUrl, nsqdTCPAddr, wsUrl)
	loader.Start()

	<-exitChan
	loader.Stop()
}
