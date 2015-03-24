package pusher

import (
	r "github.com/dancannon/gorethink"
	"log"
)

var session *r.Session

func init() {
	var err error
	session, err = r.Connect(r.ConnectOpts{
		Address:  "192.168.33.10:28015",
		Database: "mercury",
	})

	if err != nil {
		log.Fatalln(err.Error())
	}
}
