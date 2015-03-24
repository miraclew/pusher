package restful

import (
	"encoding/json"
	"net/http"
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

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ApiController struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Id             int
	Response       *ApiResponse
}

func (this *ApiController) Init(w http.ResponseWriter, r *http.Request) {
	this.ResponseWriter = w
	this.Request = r
	this.Response = &ApiResponse{Code: 0, Message: "", Data: nil}
	this.Id, _ = strconv.Atoi(this.Request.URL.Query().Get("id"))
}

func (this *ApiController) Before() error {
	return nil
}

func (this *ApiController) After() {
	if this.Response != nil {
		var response []byte
		response, _ = json.Marshal(this.Response)
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

func (this *ApiController) Ok(data interface{}) {
	this.Response.Code = 0
	this.Response.Data = data
}

func (this *ApiController) Fail(code int, data interface{}) {
	this.Response.Code = code
	this.Response.Data = data
}
