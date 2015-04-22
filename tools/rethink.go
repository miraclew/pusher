package main

import (
	"coding.net/miraclew/pusher/pusher"
	r "github.com/dancannon/gorethink"
	"log"
	"time"
)

var rdb *r.Session

type M struct {
	Id        string
	Type      int
	Text      string
	CreatedAt time.Time
}

func ttt() {
	var err error
	rdb, err = r.Connect(r.ConnectOpts{
		Address:  "192.168.33.10:28015",
		Database: "mercury",
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	// m := &M{
	// 	Type:      100,
	// 	Text:      "dddd",
	// 	CreatedAt: time.Now(),
	// }

	// res, err := r.Db("mercury").Table("messages").Insert(m).RunWrite(rdb)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// m.Id = res.GeneratedKeys[0]
	// //return m, nil
	// log.Println(m)

	res2, err2 := r.Table("messages").Get("161c9713-6795-4193-aba6-e271e63a5843").Run(rdb)
	if err2 != nil || res2.IsNil() {
		log.Fatalln(err2)
	}

	m2 := &pusher.Message{}
	res2.One(m2)

	log.Printf("%#v \n", m2.Opts["apn_nable"])
}
