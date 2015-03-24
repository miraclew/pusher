package api

import (
	"coding.net/miraclew/pusher/restful"
)

type ChannelMsgController struct {
	restful.ApiController
}

func (this *ChannelMsgController) Get() {
	this.Response.Data = "ChannelMsgController"
}
