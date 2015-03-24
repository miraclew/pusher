package restful

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

type MockResponseWriter struct {
	header http.Header
	body   bytes.Buffer
	Status int
}

func (m *MockResponseWriter) Header() http.Header {
	if m.header == nil {
		m.header = make(http.Header)
	}
	return m.header
}

func (m *MockResponseWriter) Write(b []byte) (int, error) {
	return m.body.Write(b)
}

func (m *MockResponseWriter) WriteHeader(status int) {
	m.Status = status
}

func TestApiController(t *testing.T) {
	c := &ApiController{}
	w := &MockResponseWriter{}
	r := &http.Request{}

	c.Init(w, r)
	err := c.Before()
	if err != nil {
		t.FailNow()
	}
	c.Get()
	c.After()

	if w.Status != http.StatusMethodNotAllowed {
		t.FailNow()
	}
}

type XController struct {
	ApiController
}

func (x *XController) Get() {
	x.Response.Data = "Yes, I'm the data"
}

func (x *XController) Post() {
	x.Response.Data = struct {
		Name string
		Age  int
	}{
		"Hello",
		12,
	}
}

func TestXControllerGet(t *testing.T) {
	c := &XController{}
	w := &MockResponseWriter{}
	r := &http.Request{}

	c.Init(w, r)
	err := c.Before()
	if err != nil {
		t.FailNow()
	}
	c.Get()
	c.After()

	if w.Status != http.StatusOK {
		t.Errorf("status is %d", w.Status)
	}
}

func TestXControllerPost(t *testing.T) {
	c := &XController{}
	w := &MockResponseWriter{}
	r := &http.Request{}

	c.Init(w, r)
	err := c.Before()
	if err != nil {
		t.FailNow()
	}
	c.Post()
	c.After()

	if w.Status != http.StatusOK {
		t.Errorf("status is %d", w.Status)
	}

	// var b1, b2 []byte
	// b1, _ = json.Marshal(c.Response.Data)
	// b2 = w.body.Bytes()
	// if string(b1) != string(b2) {
	// 	t.Error("data incorrect")
	// }
}

func TestTypeConv(t *testing.T) {
	md := map[string]string{
		"name": "aob",
	}
	s, err := json.Marshal(md)
	t.Log(string(s))
	s, err = json.Marshal(struct{ Name string }{"bob"})
	if err != nil {
		t.Error("Marshal failed")
	}
	t.Log(string(s))
}

type Hello interface {
	Greeting()
	What()
}

type A struct{}

func (a *A) Greeting() {

}

func (a *A) What() {

}

type B struct {
	A
}

func callWhat(h Hello) {
	h.What()
}

func TestInterface(t *testing.T) {
	//var x Hello
	callWhat(new(A))
	callWhat(new(B))
}
