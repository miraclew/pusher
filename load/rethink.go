package main

import (
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

func main() {
	var err error
	rdb, err = r.Connect(r.ConnectOpts{
		Address:  "127.0.0.1:28015",
		Database: "mercury",
	})

	if err != nil {
		log.Fatalln(err.Error())
	}

	m := &M{
		Type:      100,
		Text:      "dddd",
		CreatedAt: time.Now(),
	}

	res, err := r.Db("mercury").Table("messages").Insert(m).RunWrite(rdb)
	if err != nil {
		log.Fatalln(err)
	}

	m.Id = res.GeneratedKeys[0]
	//return m, nil
	log.Println(m)

	res2, err2 := r.Db("mercury").Table("messages").Get(m.Id).Run(rdb)
	if err2 != nil || res2.IsNil() {
		log.Fatalln(err2)
	}

	m2 := &M{}
	res2.One(m2)

	log.Println(m2)
}
