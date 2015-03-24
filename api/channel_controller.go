package api

import (
	"coding.net/miraclew/pusher/pusher"
	"coding.net/miraclew/pusher/restful"
	"github.com/miraclew/mrs/util"
)

type ChannelController struct {
	restful.ApiController
}

func (this *ChannelController) Post() {
	values := this.Request.PostForm

	members := util.SplitUniqSort(values.Get("members"))

	this.Response.Data, _ = pusher.GetChannelByMembers(members)
}
