package push

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"testing"
)

func setup() {
	var err error
	db, err = sqlx.Connect("mysql", "root:abc123@tcp(ubuntu:3306)/pusher?charset=utf8")
	if err != nil {
		log.Println(err.Error())
		return
	} else {
		log.Println("db ok")
	}

	SetDb(db)
}

func TestSave(t *testing.T) {
	setup()
	msg := NewMessage(MSG_TYPE_ACK, 123, "111", "payload", "opts")
	err := msg.Save()
	if err != nil {
		log.Fatalln(err.Error())
		t.Fail()
	}

	msg2, err := FindMessage(msg.Id)
	if err != nil {
		log.Fatalln(err.Error())
		t.Fail()
	}
	log.Println(msg2)
}
