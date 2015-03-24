package main

import (
	"flag"
	"fmt"
	"github.com/miraclew/mrs/util"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	showVersion  = flag.Bool("version", false, "print version string")
	reloadConfig = flag.Bool("reload", false, "reload config")
	httpAddress  = flag.String("http-address", "0.0.0.0:8080", "<addr>:<port> to listen on for HTTP clients")
	tcpAddress   = flag.String("tcp-address", "0.0.0.0:8081", "<addr>:<port> to listen on for TCP clients")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(util.Version("pusherd"))
		return
	}

	if *reloadConfig {
		fmt.Println("reloading config")
		// TODO: not implemented
	}

	httpAddr, err := net.ResolveTCPAddr("tcp", *httpAddress)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(util.Version("pusherd"))

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	options := &AppOptions{}
	app := NewApp(options)
	app.httpAddr = httpAddr

	app.Main()
	<-exitChan
	app.Exit()
}
