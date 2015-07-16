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
	showVersion = flag.Bool("version", false, "print version string")
	redisAddr   = flag.String("redisAddr", "127.0.0.1:6379", "<addr>:<port> (127.0.0.1:6379) redis address to connect")
	apnsDev     = flag.Bool("dev", false, "run on dev mode, apns push on dev env")
	err         error
)

var log *logging.Logger

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(Version("router"))
		return
	}

	log, err = xlog.Open("router")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer xlog.Close()

	log.Info(Version("router"))

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	options := &AppOptions{
		redisAddr: *redisAddr,
		apnsDev:   *apnsDev,
	}
	app := NewApp(options)

	app.Main()
	<-exitChan
	app.Exit()
}
