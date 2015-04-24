package main

import (
	"flag"
	"fmt"
	"github.com/miraclew/mrs/util"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	showVersion = flag.Bool("version", false, "print version string")
	httpAddr    = flag.String("httpAddr", "0.0.0.0:9010", "<addr>:<port> to listen on for HTTP clients")
	rethinkAddr = flag.String("rethinkAddr", "127.0.0.1:28015", "<addr>:<port> (127.0.0.1:28015) rethink address to connect")
	rethinkDb   = flag.String("rethinkDb", "", "rethink db name")
	redisAddr   = flag.String("redisAddr", "127.0.0.1:6379", "<addr>:<port> (127.0.0.1:6379) redis address to connect")
	apnsDev     = flag.Bool("dev", false, "run on dev mode, apns push on dev env")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(util.Version("pusherd"))
		return
	}

	if *rethinkDb == "" {
		fmt.Println("rethinkDb is required")
		return
	}

	log.Println(util.Version("pusherd"))

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	options := &AppOptions{
		rethinkAddr: *rethinkAddr,
		rethinkDb:   *rethinkDb,
		redisAddr:   *redisAddr,
		httpAddr:    *httpAddr,
		apnsDev:     *apnsDev,
	}
	app := NewApp(options)

	app.Main()
	<-exitChan
	app.Exit()
}

// func setupLog() {
// 	f, err := os.OpenFile("/var/log/pusher/pusherd.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		log.Fatalf("error opening file: %v", err)
// 	}
// 	defer f.Close()

// 	log.SetOutput(f)
// }
