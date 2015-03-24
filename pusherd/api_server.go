package main

import (
	"coding.net/miraclew/pusher/api"
	"coding.net/miraclew/pusher/restful"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func init() {
}

func httpServe(listener net.Listener) {
	log.Printf("HTTP: listening on %s", listener.Addr().String())

	handler := http.NewServeMux()

	handler.Handle("/channel", restful.NewRestfulApiHandler(new(api.ChannelController)))
	handler.Handle("/channel_msg", restful.NewRestfulApiHandler(new(api.ChannelMsgController)))
	handler.Handle("/private_msg", restful.NewRestfulApiHandler(new(api.PrivateMsgController)))

	handler.HandleFunc("/", root)

	server := &http.Server{
		Handler: handler,
	}

	err := server.Serve(listener)
	// theres no direct way to detect this error because it is not exposed
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		log.Printf("ERROR: http.Serve() - %s", err.Error())
	}

	log.Printf("HTTP: closing %s", listener.Addr().String())
}

func root(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "welcome to pusher server.")
}
