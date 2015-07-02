package pusher

import (
	"encoding/json"
	"log"
	// "strings"
	"testing"
)

func _TestMsgSendOpts(t *testing.T) {
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

func _TestMsgPayloadFix(t *testing.T) {
	var data = `{
	"sent_at": 1435743522,
	"body":{
		"start_time": 1435743522,
		"end_time": 1435743522
	}
}`
	payload := map[string]interface{}{}

	// d := json.NewDecoder(strings.NewReader(payload))
	// d.UseNumber()
	// err := d.Decode(&payload)

	err := json.Unmarshal([]byte(data), &payload)

	if err != nil {
		log.Println(err.Error())
		t.Fail()
	} else {
		log.Printf("%#v", payload)
	}

	// switch payload["sent_at"].(type) {
	// case float64:
	// 	payload["sent_at"] = int64(payload["sent_at"].(float64))
	// default:
	// }
	fixLongNumber(payload, "sent_at")

	body := payload["body"].(map[string]interface{})

	payload["body"] = fixLongNumber(body, "start_time")
	payload["body"] = fixLongNumber(body, "end_time")

	log.Printf("%#v", payload)
}

// func TestFindMessage2(t *testing.T) {
// 	msg, err := FindMessage("3a35fb5c-c864-4884-93bd-72ca239194b0")
// 	if err != nil {
// 		log.Println(err.Error())
// 		t.Fail()
// 	}

// 	log.Printf("%#v", msg.Payload)
// 	log.Printf("%#v", msg)
// }
