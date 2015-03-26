package api

import (
	"coding.net/miraclew/pusher/pusher"
	"coding.net/miraclew/pusher/restful"
	"encoding/json"
	"strconv"
)

type ChannelMsgController struct {
	restful.ApiController
}

func (this *ChannelMsgController) Post() {
	values := this.Request.PostForm

	typ, _ := strconv.ParseInt(values.Get("type"), 0, 64)

	var payload interface{}
	err := json.Unmarshal([]byte(values.Get("payload")), payload)
	if err != nil {
		this.Fail(403, nil)
	}

	m := pusher.NewMessage(
		values.Get("channel_id"),
		int(typ),
		payload,
		values.Get("sender_id"),
		// FIXME values.Get("options"),
		nil,
	)

	pusher.CreateMessage(m)
	pusher.GetHub().PushMsg(m)

	this.Ok(map[string]interface{}{
		"message": m,
	})
}
