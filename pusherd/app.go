package main

import (
	"coding.net/miraclew/pusher/api"
	"coding.net/miraclew/pusher/pusher"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"log"
	"net"
	"net/http"
	"sync"
)

type App struct {
	options   *AppOptions
	tcpAddr   *net.TCPAddr
	listener  net.Listener
	waitGroup sync.WaitGroup
	exitChan  chan int
	hub       *pusher.Hub
}

type AppOptions struct {
	rethinkAddr string
	rethinkDb   string
	redisAddr   string
	httpAddr    string
	apnsDev     bool
}

func NewApp(options *AppOptions) *App {
	a := &App{
		options:  options,
		exitChan: make(chan int),
		hub:      pusher.GetHub(),
	}

	return a
}

func NewAppOptions() *AppOptions {
	options := &AppOptions{}

	return options
}

func (a *App) Main() {
	pusher.Start(a.options.rethinkAddr, a.options.rethinkDb, a.options.redisAddr, a.options.apnsDev)
	a.startWS()
	a.startApi()
}

func (a *App) startWS() {
	p := pat.New()
	p.Get("/ws", WSHandler)
	p.Get("/", WSHandler)

	n := negroni.Classic()
	n.UseHandler(p)

	go func() {
		log.Println("Http ws listen ", a.options.httpAddr)
		err := http.ListenAndServe(a.options.httpAddr, n)
		if err != nil {
			log.Fatalln(err)
		}
	}()
}

func (a *App) startApi() {
	p := pat.New()
	p.Get("/about", api.HandleAbout)
	p.Get("/info", api.HandleInfo)
	p.Get("/mq", api.HandleMq)
	p.Post("/channel_msg", api.HandleChannelMsg)
	p.Post("/channel", api.HandleChannel)
	p.Post("/direct_msg", api.HandleDirectMsg)

	n := negroni.Classic()
	n.UseHandler(p)

	addr, err := net.ResolveTCPAddr("tcp", a.options.httpAddr)
	if err != nil {
		log.Fatal("httpAddr ResolveTCPAddr error: %s", err.Error())
		return
	}
	addr.Port += 1
	addr.IP = net.ParseIP("127.0.0.1")

	go func() {
		log.Println("Http api listen ", addr.String())
		err := http.ListenAndServe(addr.String(), n)
		if err != nil {
			log.Fatalln(err)
		}
	}()
}

func (a *App) Exit() {
	if a.listener != nil {
		a.listener.Close()
	}

	pusher.Stop()

	close(a.exitChan)
	a.waitGroup.Wait()
}
