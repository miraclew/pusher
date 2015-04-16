package main

import (
	"github.com/anachronistic/apns"
	"log"
	"time"
)

const (
	CERT_FILE = "/Users/aaaa/certificate/mercury/dev/cert.pem"
	KEY_FILE  = "/Users/aaaa/certificate/mercury/dev/key.unencrypted.pem"
)

func apn_push() {
	deviceToken := "f23028556c3213fb48834b845cb5e62eb29a14e89a421824f6cfdc2bf95f3384"

	payload := apns.NewPayload()
	payload.Alert = "你有一条新的消息"
	payload.Sound = "ping.aiff"
	payload.Badge = 1
	// if v, ok := msg.Opts["apn_alert"]; ok {
	// 	payload.Alert = v
	// }

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)

	t := time.Now()
	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", CERT_FILE, KEY_FILE)
	log.Println(time.Now().Sub(t))
	t = time.Now()

	resp := client.Send(pn)
	log.Println(time.Now().Sub(t))

	t = time.Now()
	resp = client.Send(pn)
	log.Println(time.Now().Sub(t))

	if !resp.Success {
		log.Println("apns err: ", resp.Error)
	} else {
		log.Println("apns success")
	}
}

func main() {
	apn_push()
}
