package hb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("Request body is %v, want %v", got, want)
	}
}

func testResourceID(t *testing.T, r *http.Request, want string) {
	id := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	if got := id; got != want {
		t.Errorf("%v ID is %v, want %v", r.URL.Path, id, want)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("Client BaseURL is %v, want %v", got, want)
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := NewClient(nil)
	inURL, outURL := "/foo", defaultBaseURL+"foo"
	inBody, outBody := &User{Name: "TestUser"}, `{"name":"TestUser"}`+"\n"

	req, _ := c.NewRequest("GET", inURL, inBody)

	// test that the endpoint URL was correctly added to the base URL
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%v) URL is %v, want %v", inURL, got, want)
	}

	// test that the body was correctly encoded to JSON
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%+v) Body is %v, want %v", inBody, got, want)
	}
}

func TestClient_Do(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		Foo string
	}

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"Foo":"bar"}`)
	})

	req, _ := client.NewRequest("GET", "/foo", nil)
	body := new(foo)
	client.Do(req, body)

	got, want := body, &foo{Foo: "bar"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Do %+v, want %+v", got, want)
	}
}

func TestClient_Do_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	_, err := client.Do(req, nil)
	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}
