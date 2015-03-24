package pusher

import (
	r "github.com/dancannon/gorethink"
	rds "github.com/fzzy/radix/redis"
	"log"
	"time"
)

var rdb *r.Session
var redis *rds.Client

func init() {
	var err error
	rdb, err = r.Connect(r.ConnectOpts{
		Address:  "192.168.33.10:28015",
		Database: "mercury",
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	redis, err = rds.DialTimeout("tcp", "192.168.33.10:6379", time.Duration(10)*time.Second)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
