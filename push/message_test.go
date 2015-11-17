package push

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"testing"
)

func setup() {
	var err error
	db, err = sqlx.Connect("mysql", "root:abc123@tcp(ubuntu:3306)/sun_push?charset=utf8")
	if err != nil {
		log.Println(err.Error())
		return
	} else {
		log.Println("db ok")
	}

	SetDb(db)
}

func _TestSave(t *testing.T) {
	setup()
	// msg := NewMessage(MSG_TYPE_ACK, 123, "111", 0, "payload", "opts")
	// err := msg.Save()
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// 	t.Fail()
	// }

	msg2, err := FindMessage(208)
	if err != nil {
		log.Fatalln(err.Error())
		t.Fail()
	}
	log.Println(msg2)

	p, err2 := msg2.GetPayload()
	log.Println(string(p), err2)
}

func TestParseOpt(t *testing.T) {
	msgJson := `{"type":1,"sub_type":1001,"chat_id":1049,"sender_id":1000,"receiver":"100063","body":"{\"mime\":\"text\",\"content\":{\"text\":\"aaa\"}}","extra":"{\"sender_name\":\"\\u5ba2\\u670d\",\"sender_avatar\":\"http:\\\/\\\/static.lover1314.me\\\/icons\\\/icon.png\",\"age\":20,\"love_status\":2,\"gender\":2,\"address\":\"\\u4e0a\\u6d77 \\u6d66\\u4e1c\\u65b0\\u533a\"}","opts":"{\"ttl\":0,\"offline_flag\":1,\"ack_flag\":\"1\",\"apn_flag\":\"1\",\"alert\":\"\\u5ba2\\u670d:aaa\"}","timestamp":1447727755671,"id":42231}`
	var v Message
	err := json.Unmarshal([]byte(msgJson), &v)
	if err != nil {
		t.Fail()
	}

	opt := v.ParseOpts()
	log.Printf("opt=%s result=%#v", v.Opts, opt)
}
