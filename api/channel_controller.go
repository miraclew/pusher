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

	// sort, uniq, trim
	members := util.SplitUniqSort(values.Get("members"))
	c, err := pusher.GetChannelByMembers(members)
	if err != nil {
		this.Fail(500, nil)
	}

	this.Ok(map[string]interface{}{
		"channel": c,
	})
}
