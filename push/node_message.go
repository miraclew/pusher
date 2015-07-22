package push

import (
	"encoding/json"
	"strconv"
)

const (
	NODE_CMD_PUSH = 1

	NODE_EVENT_ONLINE = 1
)

// router <=> connector node
type NodeCmd struct {
	Cmd  int
	Body string
}

type NodeCmdPush struct {
	MsgId int64
}

func (n *NodeCmdPush) Parse(body string) error {
	msgId, err := strconv.ParseInt(body, 10, 64)
	if err != nil {
		return err
	}

	n.MsgId = msgId
	return nil
}

func (n *NodeCmdPush) Marshal() ([]byte, error) {
	cmd := &NodeCmd{}
	cmd.Cmd = NODE_CMD_PUSH
	//cmd.Body =

	return json.Marshal(cmd)
}

type NodeEvent struct {
	Event  int
	NodeId int
	Body   string
}

type NodeEventOnline struct {
	UserId   int64
	IsOnline bool
}

func (n *NodeEventOnline) Parse(body string) error {
	evt := &NodeEventOnline{}
	err := json.Unmarshal([]byte(body), evt)
	if err != nil {
		return err
	}

	return nil
}
