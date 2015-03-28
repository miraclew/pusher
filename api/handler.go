package api

import (
	"coding.net/miraclew/pusher/pusher"
	"encoding/json"
	"github.com/miraclew/mrs/util"
	"net/http"
	"sort"
	"strconv"
)

func HandleChannel(res http.ResponseWriter, req *http.Request) {
	values := req.PostForm

	// sort, uniq, trim
	members := util.SplitUniqSort(values.Get("members"))
	c, err := pusher.GetChannelByMembers(members)
	if err != nil {
		respondFail(res, http.StatusInternalServerError, err.Error())
		return
	}

	respondOk(res, map[string]interface{}{
		"channel": c,
	})
}

func HandleChannelMsg(res http.ResponseWriter, req *http.Request) {
	values := req.PostForm

	typ, _ := strconv.ParseInt(values.Get("type"), 0, 64)

	var payload interface{}
	err := json.Unmarshal([]byte(values.Get("payload")), payload)
	if err != nil {
		respondFail(res, http.StatusForbidden, "payload malformed")
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

	respondOk(res, map[string]interface{}{
		"message": m,
	})
}

func HandlePrivateMsg(res http.ResponseWriter, req *http.Request) {
	values := req.PostForm

	var payload interface{}
	err := json.Unmarshal([]byte(values.Get("payload")), payload)
	if err != nil {
		respondFail(res, http.StatusForbidden, err.Error())
		return
	}

	senderId := values.Get("sender_id")
	receiverId := values.Get("receiver_id")
	members := sort.StringSlice{senderId, receiverId}
	members = sort.StringSlice(members)
	members.Sort()
	if senderId == "" || receiverId == "" || len(members) < 2 {
		respondFail(res, http.StatusForbidden, "required params missed")
		return
	}

	var channel *pusher.Channel
	channel, err = pusher.GetChannelByMembers(members)
	if err != nil {
		respondFail(res, http.StatusInternalServerError, err.Error())
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

	respondOk(res, map[string]interface{}{
		"message": m,
	})
}

func HandleAbout(res http.ResponseWriter, req *http.Request) {
	respondOk(res, map[string]interface{}{
		"message": "welcome to push service",
	})
}
