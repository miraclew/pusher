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
	pusher.Start(a.options.rethinkAddr, a.options.rethinkDb, a.options.redisAddr)

	p := pat.New()

	p.Get("/about", api.HandleAbout)
	p.Post("/channel_msg", api.HandleChannelMsg)
	p.Post("/channel", api.HandleChannel)
	p.Post("/direct_msg", api.HandleDirectMsg)

	p.Get("/ws", WSHandler)
	p.Get("/", WSHandler)

	n := negroni.Classic()
	n.UseHandler(p)

	go func() {
		log.Println("http listen ", a.options.httpAddr)
		err := http.ListenAndServe(a.options.httpAddr, n)
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
