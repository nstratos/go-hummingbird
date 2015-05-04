package hb

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("Client BaseURL is %v, want %v", got, want)
	}
}

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
