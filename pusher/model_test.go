package pusher

import (
	"encoding/json"
	"log"
	"testing"
)

func TestMsgSendOpts(t *testing.T) {
	opts := &MsgSendOpts{}
	err := json.Unmarshal([]byte(`{
"apnEnable": true ,
"AckEnable": false,
"apn_alert":  "你收到一条消息" ,
"offlineEnable": true ,
"ttl": 10
}`), opts)

	if err != nil {
		log.Println(err.Error())
		t.Fail()
	} else {
		log.Printf("%#v", opts)
	}
}

func TestFindMessage2(t *testing.T) {
	msg, err := FindMessage("00000820-d618-4880-aab7-bc32a756049d")
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	log.Printf("%#v", msg.Opts)
	log.Printf("%#v", msg)
}
