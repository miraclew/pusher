package main

import (
	"coding.net/miraclew/pusher/push"
	"coding.net/miraclew/pusher/xlog"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/op/go-logging"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	app              *App
	showVersion      = flag.Bool("version", false, "print version stringx")
	wsIp             = flag.String("ws-ip", "0.0.0.0", "<ip> to listen on for WebSocket clients")
	wsPort           = flag.Int("ws-port", 0, "<port> to listen on for WebSocket clients, 0 means calculate automaticly")
	nodeId           = flag.Int("node-id", 0, "id of the connector")
	clientTimeout    = flag.Int("client-timeout", 3600, "id of the connector")
	apiAddr          = flag.String("api-addr", "127.0.0.1:9011", "<addr>:<port> to listen on for Http Api clients")
	redisAddr        = flag.String("redis", "127.0.0.1:6379", "<addr>:<port> (127.0.0.1:6379) redis address to connect")
	nsqdTCPAddrs     = push.StringArray{}
	lookupdHTTPAddrs = push.StringArray{}
)

var log *logging.Logger

func init() {
	flag.Var(&nsqdTCPAddrs, "nsqd-tcp-address", "(127.0.0.1:4150) nsqd TCP address (may be given multiple times)")
	flag.Var(&lookupdHTTPAddrs, "lookupd-http-address", "(127.0.0.1:4161) lookupd HTTP address (may be given multiple times)")
	var err error
	log, err = xlog.Open("connector")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	flag.Parse()
	defer xlog.Close()

	if *showVersion {
		fmt.Println(Version("connector"))
		return
	}

	log.Info(Version("connector"))

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if *nodeId == 0 {
		fmt.Println("node-id is required")
		return
	}

	if len(nsqdTCPAddrs) == 0 && len(lookupdHTTPAddrs) == 0 {
		log.Fatalf("--nsqd-tcp-address or --lookupd-http-address required.")
	}

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var wsPortStart int64
	portStart := os.Getenv("WS_PORT_START")
	if portStart != "" {
		wsPortStart, _ = strconv.ParseInt(portStart, 10, 64)
	} else {
		wsPortStart = 9100
	}

	options := &AppOptions{
		redisAddr:        *redisAddr,
		wsIp:             *wsIp,
		wsPort:           *wsPort,
		wsPortStart:      int(wsPortStart),
		nodeId:           *nodeId,
		clientTimeout:    *clientTimeout,
		nsqdTCPAddrs:     nsqdTCPAddrs,
		lookupdHTTPAddrs: lookupdHTTPAddrs,
	}

	app = NewApp(options)

	app.Main()
	<-exitChan
	app.Exit()
}
