package api

import (
	"coding.net/miraclew/pusher/restful"
)

type ChannelController struct {
	restful.ApiController
}

func (this *ChannelController) Get() {
	this.Response.Data = "ChannelController"
}
