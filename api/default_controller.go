package api

import (
	"coding.net/miraclew/pusher/restful"
)

type DefaultController struct {
	restful.ApiController
}

func (this *DefaultController) Get() {
	this.Response.Data = "DefaultController"
}
