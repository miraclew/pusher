package api

import (
	"coding.net/miraclew/pusher/pusher"
	"encoding/json"
	"github.com/miraclew/mrs/util"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	nbApiRequests int
	startAt       time.Time
)

func init() {
	startAt = time.Now()
}

func HandleChannel(res http.ResponseWriter, req *http.Request) {
	nbApiRequests++
	data, _ := ioutil.ReadAll(req.Body)
	var body map[string]interface{}
	err := json.Unmarshal(data, &body)
	if err != nil {
		log.Printf("body malformed: body=%s err=%s", string(data), err.Error())
		respondFail(res, http.StatusBadRequest, "body malformed: "+err.Error())
		return
	}

	log.Printf("HandleChannelMsg(%#v)", body)

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
	nbApiRequests++
	data, _ := ioutil.ReadAll(req.Body)

	var msg pusher.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("body malformed: body=%s err=%s", string(data), err.Error())
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

func HandleDirectMsg(res http.ResponseWriter, req *http.Request) {
	nbApiRequests++
	data, _ := ioutil.ReadAll(req.Body)

	var msg pusher.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("body malformed: body=%s err=%s", string(data), err.Error())
		respondFail(res, http.StatusBadRequest, "body malformed: "+err.Error())
		return
	}
	//log.Printf("HandleDirectMsg(%#v)", msg)

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

func HandleAbout(res http.ResponseWriter, req *http.Request) {
	nbApiRequests++
	respondOK(res, map[string]interface{}{
		"message": "welcome to push service",
	})
}

func HandleInfo(res http.ResponseWriter, req *http.Request) {

	respondOK(res, map[string]interface{}{
		"start_at":     startAt,
		"uptime":       time.Now().Sub(startAt).Seconds(),
		"api_requests": nbApiRequests,
		"clients":      pusher.GetHub().ConnectionsCount(),
	})
}

func HandleMq(res http.ResponseWriter, req *http.Request) {
	userId := req.URL.Query().Get("user_id")
	if len(userId) <= 0 {
		respondFail(res, http.StatusBadRequest, "user_id is required")
		return
	}

	ids, messages, err := pusher.GetUserQueuedMessages(userId)
	if err != nil {
		respondOK(res, map[string]interface{}{
			"err": err.Error(),
		})
	}
	respondOK(res, map[string]interface{}{
		"count":    len(ids),
		"ids":      ids,
		"messages": messages,
	})
}
