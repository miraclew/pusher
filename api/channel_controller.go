package api

import (
	"coding.net/miraclew/pusher/restful"
	"crypto/md5"
	"encoding/hex"
	r "github.com/dancannon/gorethink"
	"github.com/miraclew/mrs/util"
	"io"
	"log"
	"strings"
)

var session *r.Session

func init() {
	var err error
	session, err = r.Connect(r.ConnectOpts{
		Address:  "l5.local:28015",
		Database: "mercury",
	})

	if err != nil {
		log.Fatalln(err.Error())
	}
}

type ChannelController struct {
	restful.ApiController
}

func (this *ChannelController) Post() {
	values := this.Request.PostForm

	members := util.SplitUniqSort(values.Get("members"))
	h := md5.New()
	io.WriteString(h, strings.Join(members, ","))
	hash := hex.EncodeToString(h.Sum(nil))

	query := r.Db("mercury").Table("channels").Filter(r.Row.Field("hash"), hash).Limit(1)
	res, _ := query.Run(session)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var x interface{}
	res.One(&x)

	this.Response.Data = x
}
