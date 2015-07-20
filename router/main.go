package main

import (
	"coding.net/miraclew/pusher/app"
	"coding.net/miraclew/pusher/xlog"
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"os"
	"os/signal"
	"syscall"
)

var (
	showVersion      = flag.Bool("version", false, "print version string")
	redisAddr        = flag.String("redisAddr", "127.0.0.1:6379", "<addr>:<port> (127.0.0.1:6379) redis address to connect")
	apnsDev          = flag.Bool("dev", false, "run on dev mode, apns push on dev env")
	nsqdTCPAddrs     = app.StringArray{}
	lookupdHTTPAddrs = app.StringArray{}
)

var log *logging.Logger

func init() {
	flag.Var(&nsqdTCPAddrs, "nsqd-tcp-address", "nsqd TCP address (may be given multiple times)")
	flag.Var(&lookupdHTTPAddrs, "lookupd-http-address", "lookupd HTTP address (may be given multiple times)")
	var err error
	log, err = xlog.Open("router")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	flag.Parse()
	defer xlog.Close()

	if *showVersion {
		fmt.Println(Version("router"))
		return
	}

	if len(nsqdTCPAddrs) == 0 && len(lookupdHTTPAddrs) == 0 {
		log.Fatalf("--nsqd-tcp-address or --lookupd-http-address required.")
	}

	if len(nsqdTCPAddrs) != 0 && len(lookupdHTTPAddrs) != 0 {
		log.Fatalf("use --nsqd-tcp-address or --lookupd-http-address not both")
	}

	log.Info(Version("router"))

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	options := &AppOptions{
		redisAddr:        *redisAddr,
		nsqdTCPAddrs:     nsqdTCPAddrs,
		lookupdHTTPAddrs: lookupdHTTPAddrs,
		apnsDev:          *apnsDev,
	}
	app := NewApp(options)

	app.Main()
	<-exitChan
	app.Exit()
}
