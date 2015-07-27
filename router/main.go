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
	showVersion = flag.Bool("version", false, "print version string")
	redisAddr   = flag.String("redis", "127.0.0.1:6379", "<addr>:<port> (127.0.0.1:6379) redis address to connect")
	// mysqlAddr        = flag.String("mysql", "", "user:pass@tcp(localhost:3306)/pusher?charset=utf8 mysql address to connect")
	apnsDev          = flag.Bool("dev", false, "run on dev mode, apns push on dev env")
	nsqdTCPAddrs     = push.StringArray{}
	lookupdHTTPAddrs = push.StringArray{}
	app              *App
)

var log *logging.Logger

func init() {
	flag.Var(&nsqdTCPAddrs, "nsqd-tcp-address", "(127.0.0.1:4150) nsqd TCP address (may be given multiple times)")
	flag.Var(&lookupdHTTPAddrs, "lookupd-http-address", "(127.0.0.1:4161) lookupd HTTP address (may be given multiple times)")
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

	if len(nsqdTCPAddrs) == 0 || len(lookupdHTTPAddrs) == 0 {
		log.Fatalf("--nsqd-tcp-address and --lookupd-http-address required.")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mysqlAddr := os.Getenv("MYSQL_DSN")
	return

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
		mysqlAddr:        mysqlAddr,
		apnsDev:          *apnsDev,
	}
	app = NewApp(options)

	app.Main()
	<-exitChan
	app.Exit()
}
