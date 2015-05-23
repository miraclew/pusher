package pusher

import (
	r "github.com/dancannon/gorethink"
	"github.com/garyburd/redigo/redis"
	"log"
)

const (
	OPT_RETHINK_ADDRESS = 1
	OPT_RETHINK_DB      = 2
	OPT_REDIS_ADDRESS   = 3
)

var (
	pool       *redis.Pool
	rdb        *r.Session
	apnsDev    bool
	_redisAddr string
)

func init() {
}

func Start(rethinkAddress string, rethinkDb string, redisAddr string, devMode bool) {
	_redisAddr = redisAddr
	var err error
	rdb, err = r.Connect(r.ConnectOpts{
		Address:  rethinkAddress,
		Database: rethinkDb,
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("rethink db connected: %s/%s", rethinkAddress, rethinkDb)

	pool = newRedisPool(redisAddr, "")
	apnsDev = devMode
}

func Stop() {
	rdb.Close()
	pool.Close()
}
