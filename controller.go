package restful

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
)

type Controller interface {
	// Request lifetime
	Init(w http.ResponseWriter, r *http.Request)
	Before() error
	After()
	// Http methods
	Get()    // Retrieving Objects OR Queries
	Post()   // Creating Objects
	Put()    // Updating Object
	Delete() // Delete Object
	Head()
	Patch()
	Options()
}

type ApiController struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Id             int
	Data           interface{}
}

func (this *ApiController) Init(w http.ResponseWriter, r *http.Request) {
	this.ResponseWriter = w
	this.Request = r
	this.Id, _ = strconv.Atoi(this.Request.URL.Query().Get("id"))
}

func (this *ApiController) Before() error {
	return nil
}

func (this *ApiController) After() {
	if this.Data != nil {
		var response []byte
		k := reflect.TypeOf(this.Data).Kind()
		if k == reflect.String {
			response = []byte(this.Data.(string))
		} else {
			response, _ = json.Marshal(this.Data)
		}

		this.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		//this.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(response)))
		this.ResponseWriter.WriteHeader(http.StatusOK)
		this.ResponseWriter.Write(response)
	}
}

func (this *ApiController) Get() {
	http.Error(this.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (this *ApiController) Post() {
	http.Error(this.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (this *ApiController) Put() {
	http.Error(this.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (this *ApiController) Delete() {
	http.Error(this.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (this *ApiController) Head() {
	http.Error(this.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (this *ApiController) Patch() {
	http.Error(this.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (this *ApiController) Options() {
	http.Error(this.ResponseWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
