package pusher

import (
	r "github.com/dancannon/gorethink"
	rds "github.com/fzzy/radix/redis"
	"log"
	"time"
)

const (
	OPT_RETHINK_ADDRESS = 1
	OPT_RETHINK_DB      = 2
	OPT_REDIS_ADDRESS   = 3
)

var rdb *r.Session
var redis *rds.Client

func init() {
}

func Start(rethinkAddress string, rethinkDb string, redisAddress string) {
	var err error
	rdb, err = r.Connect(r.ConnectOpts{
		Address:  rethinkAddress,
		Database: rethinkDb,
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("rethink db connected: %s/%s", rethinkAddress, rethinkDb)

	redis, err = rds.DialTimeout("tcp", redisAddress, time.Duration(10)*time.Second)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("redis connected: %s", redisAddress)
}

func Stop() {
	rdb.Close()
	redis.Close()
}
