package main

import (
	"github.com/anachronistic/apns"
	"log"
)

const (
	CERT_FILE = "/Users/aaaa/certificate/mercury/dev/cert.pem"
	KEY_FILE  = "/Users/aaaa/certificate/mercury/dev/key.unencrypted.pem"
)

func main() {
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

	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", CERT_FILE, KEY_FILE)
	resp := client.Send(pn)

	if !resp.Success {
		log.Println("apns err: ", resp.Error)
	} else {
		log.Println("apns success")
	}
}
