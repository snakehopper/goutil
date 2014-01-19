package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

type TestStruct struct {
	Id int
}

func TestReturnJson(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		s := TestStruct{Id: 1}
		ReturnJson(w, s)
	})

	resp, _ := http.Get(server.URL)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	want := []byte("{\n\t\"Id\": 1\n}")
	if !reflect.DeepEqual(body, want) {
		t.Errorf("ReturnJson returned %+v, want %+v", body, want)
	}
}

func TestGetUrl(t *testing.T) {
	setup()
	defer teardown()

	want := `{"id":1}`
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, want)
	})

	res := GetUrl(nil, server.URL)

	if !reflect.DeepEqual(res, want) {
		t.Errorf("GetUrl returned %+v, want %+v", res, want)
	}
}
