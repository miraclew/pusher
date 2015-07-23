package push

import ()

const (
	NODE_CMD_PUSH = 1

	NODE_EVENT_ONLINE = 1
)

// router <=> connector node
type NodeCmd struct {
	Cmd  int
	Body []byte
}

type NodeCmdPush struct {
	MsgId      int64
	ReceiverId int64
	Payload    []byte
}

type NodeEvent struct {
	Event  int
	NodeId int
	Body   []byte
}

type NodeEventOnline struct {
	UserId   int64
	IsOnline bool
}
