package api

import (
	"coding.net/miraclew/pusher/restful"
)

type PrivateMsgController struct {
	restful.ApiController
}

func (this *PrivateMsgController) Get() {
	this.Response.Data = "PrivateMsgController"
}
