package api

import (
	"coding.net/miraclew/pusher/pusher"
	"coding.net/miraclew/pusher/restful"
	"encoding/json"
	"sort"
	// "strconv"
)

type PrivateMsgController struct {
	restful.ApiController
}

func (this *PrivateMsgController) Post() {
	values := this.Request.PostForm

	//typ, _ := strconv.ParseInt(values.Get("type"), 0, 64)

	var payload interface{}
	err := json.Unmarshal([]byte(values.Get("payload")), payload)
	if err != nil {
		this.Fail(403, nil)
		return
	}

	senderId := values.Get("sender_id")
	receiverId := values.Get("receiver_id")
	members := sort.StringSlice{senderId, receiverId}
	members = sort.StringSlice(members)
	members.Sort()
	if senderId == "" || receiverId == "" || len(members) < 2 {
		this.Fail(ERR_INVALID_ARGS, nil)
		return
	}

	var channel *pusher.Channel
	channel, err = pusher.GetChannelByMembers(members)
	if err != nil {
		this.Fail(ERR_INTERAL_ERROR, nil)
		return
	}

	options := map[string]interface{}{
		"receiver_id": receiverId,
	}
	m := pusher.NewMessage(
		channel.Id,
		pusher.MESSAGE_TYPE_NORMAL,
		payload,
		senderId,
		options,
	)

	pusher.CreateMessage(m)
	pusher.GetHub().PushMsg(m)

	this.Ok(map[string]interface{}{
		"message": m,
	})
}
