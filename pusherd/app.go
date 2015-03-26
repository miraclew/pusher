package main

import (
	"coding.net/miraclew/pusher/pusher"
	"log"
	"net"
	"sync"
)

type App struct {
	options      *AppOptions
	tcpAddr      *net.TCPAddr
	httpAddr     *net.TCPAddr
	httpListener net.Listener
	waitGroup    sync.WaitGroup
	exitChan     chan int
	hub          *pusher.Hub
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

	httpListener, err := net.Listen("tcp", a.httpAddr.String())
	if err != nil {
		log.Fatalf("FATAL: listen (%s) failed - %s", a.tcpAddr, err.Error())
	}
	a.httpListener = httpListener

	a.waitGroup.Add(2)
	go func() {
		httpServe(httpListener)
		a.waitGroup.Done()
	}()

	go func() {
		wsServe(httpListener, a.hub)
		a.waitGroup.Done()
	}()
}

func (a *App) Exit() {
	if a.httpListener != nil {
		a.httpListener.Close()
	}

	pusher.Stop()

	close(a.exitChan)
	a.waitGroup.Wait()
}
