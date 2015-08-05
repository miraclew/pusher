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
	"syscall"
)

var (
	app              *App
	showVersion      = flag.Bool("version", false, "print version string")
	sandbox          = flag.Bool("sandbox", false, "connect to sandbox server")
	nsqdTCPAddrs     = push.StringArray{}
	lookupdHTTPAddrs = push.StringArray{}
)

var log *logging.Logger

func init() {
	flag.Var(&nsqdTCPAddrs, "nsqd-tcp-address", "(127.0.0.1:4150) nsqd TCP address (may be given multiple times)")
	flag.Var(&lookupdHTTPAddrs, "lookupd-http-address", "(127.0.0.1:4161) lookupd HTTP address (may be given multiple times)")
	var err error
	log, err = xlog.Open("apns")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	flag.Parse()
	defer xlog.Close()
	godotenv.Load()

	if *showVersion {
		fmt.Println(Version("apns"))
		return
	}

	log.Info(Version("apns"))

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

	options := &AppOptions{
		nsqdTCPAddrs:     nsqdTCPAddrs,
		lookupdHTTPAddrs: lookupdHTTPAddrs,
		sandboxCert:      os.Getenv("SANDBOX_CERT"),
		sandboxKey:       os.Getenv("SANDBOX_KEY"),
		prodCert:         os.Getenv("PROD_CERT"),
		prodKey:          os.Getenv("PROD_KEY"),
	}

	app = NewApp(options)

	app.Main()
	<-exitChan
	app.Exit()
}
