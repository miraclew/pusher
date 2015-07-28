package main

import (
	"coding.net/miraclew/pusher/push"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

type Router struct {
	producers map[string]*nsq.Producer
}

func NewRouter() *Router {
	producers := make(map[string]*nsq.Producer)
	cfg := nsq.NewConfig()
	// cfg.UserAgent = fmt.Sprintf("to_nsq/%s go-nsq/%s", version.Binary, nsq.VERSION)

	for _, addr := range app.options.nsqdTCPAddrs {
		producer, err := nsq.NewProducer(addr, cfg)
		if err != nil {
			log.Fatalf("failed to create nsq.Producer - %s", err)
		}
		producers[addr] = producer
	}

	if len(producers) == 0 {
		log.Fatal("--nsqd-tcp-address required")
	}

	return &Router{
		producers: producers,
	}
}

func (r *Router) route(msg *push.Message) error {
	receivers, err := msg.ParseReceivers()
	if err != nil {
		log.Error("Bad Receivers: %s", err.Error())
	}

	for _, receiver := range receivers {
		err := r.routeDirect(receiver, msg)
		if err != nil {

		}
	}

	return nil
}

func (r *Router) routeDirect(userId int64, msg *push.Message) error {
	log.Info("routeDirect(%d, %d)", userId, msg.Id)
	client, err := push.GetClient(userId)
	if err != nil {
		log.Error("GetClient %d error: %s", userId, err.Error())
		return err
	}

	apnFlag := msg.ParseOpts().ApnFlag
	if apnFlag == push.MSG_OPT_APN_NOTIFY_ONLY && client.DeviceType == push.DEVICE_TYPE_IOS {
		return r.publishToApns(userId, msg)
	}

	if msg.ParseOpts().OfflineFlag != push.MSG_OPT_OFFLINE_DISABLE {
		_, err = r.pushToQueue(userId, msg.Id)
		if err != nil {
			log.Error("pushToQueue err=%s", err.Error())
			return err
		}
	}

	if client.IsOnline() {
		err := r.publishToNode(client.NodeId, userId, msg)
		if err != nil {
			log.Error("publishToNode(%d, %d, %#v) error: %s", client.NodeId, userId, msg, err.Error())
			return err
		}
	} else {
		if client.DeviceType == push.DEVICE_TYPE_IOS && msg.ParseOpts().ApnFlag == push.MSG_OPT_APN_DEFAULT {
			log.Info("client is offline: %d, pushToIosDevice", client.UserId)
			return r.publishToApns(userId, msg)
		} else {
			log.Info("client is offline: %d, skip", client.UserId)
		}
	}

	return nil
}

func (r *Router) pushToQueue(userId int64, msgId int64) (int, error) {
	log.Info("pushToQueue(%d, %d)", userId, msgId)
	conn := app.redisPool.Get()
	defer conn.Close()

	cmd := "rpush"
	return redis.Int(conn.Do(cmd, fmt.Sprintf("mq:%d", userId), msgId))
}

func (r *Router) publishToNode(nodeId int, userId int64, msg *push.Message) error {
	log.Info("publishToNode(%d, %d) to %d producers", nodeId, msg.Id, len(r.producers))
	cmd := &push.NodeCmd{}
	cmd.Cmd = push.NODE_CMD_PUSH

	var err error
	body := &push.NodeCmdPush{}
	body.MsgId = msg.Id
	body.Payload, err = msg.GetPayload()
	if err != nil {
		return err
	}

	body.ReceiverId = userId
	cmd.Body, err = json.Marshal(body)
	if err != nil {
		return err
	}

	b, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	for _, producer := range r.producers {
		err := producer.Publish(fmt.Sprintf("connector-%d", nodeId), b)
		if err != nil {
			log.Error("Publish error: %s", err.Error())
		}
	}

	return nil
}

func (r *Router) processQueue(userId int64) error {
	log.Info("processQueue %d", userId)
	client, err := push.GetClient(userId)
	if err != nil {
		log.Error("processQueue GetClient %d error: %s", userId, err.Error())
		return err
	}
	if !client.IsOnline() {
		log.Warning("Client %d is not online when processQueue", userId)
		return nil
	}

	conn := app.redisPool.Get()
	defer conn.Close()
	key := fmt.Sprintf("mq:%d", userId)
	ids, err := redis.Strings(conn.Do("lrange", key, 0, -1))
	if err != nil {
		return err
	}

	log.Info("queue(%s) size: %d", key, len(ids))
	if len(ids) <= 0 {
		return nil
	}

	for i := 0; i < len(ids); i++ {
		var msgId int64
		msgId, err = strconv.ParseInt(ids[i], 10, 64)
		if err != nil {
		}
		var msg *push.Message
		msg, err = push.FindMessage(msgId)
		if err != nil {
			log.Error(fmt.Sprintf("FindMessage(%s): %s", msgId, err.Error()))
			continue
		}
		if msg == nil {
			log.Error(fmt.Sprintf("msgId: %s not found", msgId))
			_, err := redis.Int(conn.Do("lrem", key, 0, msgId))
			if err != nil {
				log.Error("lrem: %s \n", err.Error())
			}
			continue
		}

		err = r.publishToNode(client.NodeId, userId, msg)
		if err != nil {
			log.Error("publishToNode(%d, %d, %#v) error: %s", client.NodeId, userId, msg, err.Error())
			continue
		}
	}
	return nil
}

func (r *Router) publishToApns(userId int64, msg *push.Message) error {
	conn := app.redisPool.Get()
	defer conn.Close()

	deviceToken, err := redis.String(conn.Do("get", fmt.Sprintf("apn_u2t:%d", userId)))
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			log.Warning("pushToIosDevice %d, deviceToken not exist", userId)
			return nil
		} else {
			log.Error("pushToIosDevice error: %s", err.Error())
			return err
		}
	}

	if deviceToken == "" {
		log.Warning("pushToIosDevice %d deviceToken is empty", userId)
		return nil
	}

	log.Info("publishToApns(%d, %d) to %d producers", userId, msg.Id, len(r.producers))
	cmd := &push.ApnsCmd{}
	cmd.DeviceToken = deviceToken
	cmd.MsgId = fmt.Sprintf("%d", msg.Id)
	cmd.UserId = userId
	cmd.Alert = msg.ParseOpts().Alert

	b, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	for _, producer := range r.producers {
		err := producer.Publish("apns", b)
		if err != nil {
			log.Error("Publish error: %s", err.Error())
		}
	}

	return nil
}

func (r *Router) routeChannel(channelId string, msg *push.Message) error {
	return nil
}
