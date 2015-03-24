package restful

import (
	"net/http"
)

type RestfulApiHandler struct {
	controller Controller
}

func (this *RestfulApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.controller.Init(w, r)
	err := this.controller.Before()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if r.Method == "GET" {
		this.controller.Get()
	} else if r.Method == "POST" {
		r.ParseForm()
		this.controller.Post()
	} else if r.Method == "PUT" {
		r.ParseForm()
		this.controller.Put()
	} else if r.Method == "DELETE" {
		this.controller.Delete()
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	this.controller.After()
}

func NewRestfulApiHandler(c Controller) *RestfulApiHandler {
	return &RestfulApiHandler{controller: c}
}
