package main

import (
	"coding.net/miraclew/pusher/xlog"
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"os"
	"os/signal"
	"syscall"
)

var (
	log         *logging.Logger
	showVersion = flag.Bool("version", false, "print version string")
	wsIp        = flag.String("wsIp", "0.0.0.0", "<ip> to listen on for WebSocket clients")
	wsPort      = flag.Int("wsPort", 9100, "<port> to listen on for WebSocket clients")
	nodeId      = flag.Int("nodeId", 1, "id of the connector")
	apiAddr     = flag.String("apiAddr", "127.0.0.1:9011", "<addr>:<port> to listen on for Http Api clients")
	redisAddr   = flag.String("redisAddr", "127.0.0.1:6379", "<addr>:<port> (127.0.0.1:6379) redis address to connect")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(Version("connector"))
		return
	}

	var err error
	log, err = xlog.Open("connector")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer xlog.Close()

	log.Info(Version("connector"))

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	options := &AppOptions{
		redisAddr: *redisAddr,
		wsIp:      *wsIp,
		wsPort:    *wsPort,
		nodeId:    *nodeId,
	}

	app := NewApp(options)

	app.Main()
	<-exitChan
	app.Exit()
}
