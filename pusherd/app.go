package main

import (
	"coding.net/miraclew/pusher/api"
	"coding.net/miraclew/pusher/pusher"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"net"
	"net/http"
	"sync"
)

type App struct {
	options   *AppOptions
	tcpAddr   *net.TCPAddr
	httpAddr  *net.TCPAddr
	listener  net.Listener
	waitGroup sync.WaitGroup
	exitChan  chan int
	hub       *pusher.Hub
}

type AppOptions struct {
	rethinkAddr string
	rethinkDb   string
	redisAddr   string
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
	p.Get("/ws", WSHandler)

	p.Get("/", WSHandler)
	p.Get("/about", api.HandleAbout)
	p.Post("/channel_msg", api.HandleChannelMsg)
	p.Post("/channel", api.HandleChannel)
	p.Post("/private_msg", api.HandlePrivateMsg)

	n := negroni.Classic()
	n.UseHandler(p)

	go func() {
		http.ListenAndServe(":9010", n)
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
