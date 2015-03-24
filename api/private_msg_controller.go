package api

import (
	"coding.net/miraclew/pusher/pusher"
	"coding.net/miraclew/pusher/restful"
	"encoding/json"
	"github.com/miraclew/mrs/util"
	"strconv"
)

type PrivateMsgController struct {
	restful.ApiController
}

func (this *PrivateMsgController) Post() {
	values := this.Request.PostForm

	typ, _ := strconv.ParseInt(values.Get("type"), 0, 64)

	var payload interface{}
	err := json.Unmarshal([]byte(values.Get("payload")), payload)
	if err != nil {
		return this.Fail(403, nil)
	}

	senderId := values.Get("sender_id")
	receiverId := values.Get("receiver_id")
	members := util.SplitUniqSort([]string{senderId, receiverId})
	if senderId == nil || receiverId == nil || len(members) < 2 {
		return this.Fail(ERR_INVALID_ARGS, nil)
	}

	var channel *pusher.Channel
	channel, err = pusher.GetChannelByMembers(members)
	if err != nil {
		return this.Fail(ERR_INTERAL_ERROR, nil)
	}

	options := map[string]interface{}{
		"receiver_id": receiverId,
	}
	m := pusher.NewMessage(
		channel.Id,
		MESSAGE_TYPE_NORMAL,
		payload,
		senderId,
		options,
	)

	pusher.CreateMessage(m)
	// TODO: PUSH THE MESSAGE

	this.Ok(map[string]interface{}{
		"message": m,
	})
}
