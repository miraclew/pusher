package api

import (
	"coding.net/miraclew/pusher/pusher"
	"encoding/json"
	"github.com/miraclew/mrs/util"
	"log"
	"net/http"
	"sort"
	"strconv"
)

func HandleChannel(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	log.Printf("HandleChannel(%#v)", req.Form)

	// sort, uniq, trim
	members := util.SplitUniqSort(req.PostFormValue("members"))
	if len(members) <= 0 {
		respondFail(res, http.StatusBadRequest, "members is required")
		return
	}

	c, err := pusher.GetChannelByMembers(members)
	if err != nil {
		respondFail(res, http.StatusInternalServerError, err.Error())
		return
	}

	respondOK(res, map[string]interface{}{
		"channel": c,
	})
}

func HandleChannelMsg(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	values := req.PostForm

	log.Printf("HandleChannelMsg(%v)", values)
	typ, _ := strconv.ParseInt(values.Get("type"), 0, 64)

	channelId := values.Get("channel_id")
	if len(channelId) <= 0 {
		respondFail(res, http.StatusBadRequest, "channel_id is required")
		return
	}

	senderId := values.Get("sender_id")
	if len(senderId) <= 0 {
		respondFail(res, http.StatusBadRequest, "sender_id is required")
		return
	}

	var payload interface{}
	err := json.Unmarshal([]byte(values.Get("payload")), payload)
	if err != nil {
		respondFail(res, http.StatusBadRequest, "payload malformed")
		return
	}

	m := pusher.NewMessage(
		channelId,
		int(typ),
		payload,
		senderId,
		// FIXME values.Get("options"),
		nil,
	)

	pusher.CreateMessage(m)
	pusher.GetHub().PushMsg(m)

	respondOK(res, map[string]interface{}{
		"message": m,
	})
}

func HandlePrivateMsg(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	values := req.PostForm

	log.Printf("HandlePrivateMsg(%v)", values)
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

	respondOK(res, map[string]interface{}{
		"message": m,
	})
}

func HandleAbout(res http.ResponseWriter, req *http.Request) {
	respondOK(res, map[string]interface{}{
		"message": "welcome to push service",
	})
}
