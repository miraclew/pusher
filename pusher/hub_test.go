package pusher

import (
	"fmt"
	"testing"
)

func init() {
	//Start("192.168.33.10:28015", "mercury", "192.168.33.10:6379", true)
}

func _TestApns(t *testing.T) {
	hub := GetHub()
	opts := &MsgSendOpts{}
	opts.Alert = "hello apns"

	msg := NewMessage("aaa", 1, map[string]interface{}{"name": "hello"}, "1000", opts)
	err := hub.pushToIosDevice(10145, msg, 1)

	if err != nil {
		t.Fail()
		fmt.Println(err.Error())
	}
}
