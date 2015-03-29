package api

import (
	"coding.net/miraclew/pusher/pusher"
	"encoding/json"
	// "fmt"
	"github.com/miraclew/mrs/util"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	// "strconv"
)

func HandleChannel(res http.ResponseWriter, req *http.Request) {
	data, _ := ioutil.ReadAll(req.Body)
	var body map[string]interface{}
	err := json.Unmarshal(data, &body)
	if err != nil {
		log.Println("body malformed: " + err.Error())
		respondFail(res, http.StatusBadRequest, "body malformed: "+err.Error())
		return
	}

	// sort, uniq, trim
	members := util.SplitUniqSort(body["members"].(string))
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
	data, _ := ioutil.ReadAll(req.Body)

	var msg pusher.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println("body malformed: " + err.Error())
		respondFail(res, http.StatusBadRequest, "body malformed: "+err.Error())
		return
	}
	log.Printf("HandleChannelMsg(%#v)", msg)

	if len(msg.ChannelId) <= 0 {
		respondFail(res, http.StatusBadRequest, "channel_id is required")
		return
	}

	if len(msg.SenderId) <= 0 {
		respondFail(res, http.StatusBadRequest, "sender_id is required")
		return
	}

	if msg.Payload == nil {
		respondFail(res, http.StatusBadRequest, "payload is required")
		return
	}

	pusher.CreateMessage(&msg)
	pusher.GetHub().PushMsg(&msg)

	respondOK(res, map[string]interface{}{
		"message": msg,
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
