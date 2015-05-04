package hb

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	client *Client

	server *httptest.Server

	mux *http.ServeMux
)

func setup() {

	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(nil)
	client.BaseURL, _ = url.Parse(server.URL)
}

func teardown() {
	server.Close()
}

func TestUserService_Authenticate(t *testing.T) {
	setup()
	defer teardown()

	client.User.SetCredentials("TestUser", "", "TestPass")

	mux.HandleFunc("/users/authenticate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"username":"TestUser","email":"","password":"TestPass"}`)
		fmt.Fprintf(w, `"token1234"`)
	})

	token, err := client.User.Authenticate()
	if err != nil {
		t.Errorf("User.Authenticate returned error %v", err)
	}
	if got, want := token, "token1234"; got != want {
		t.Errorf("User.Authenticate token is %v, want %v", got, want)
	}
}

func TestUserService_Authenticate_unauthorized(t *testing.T) {
	setup()
	defer teardown()

	client.User.SetCredentials("InvalidTestUser", "", "TestPass")

	mux.HandleFunc("/users/authenticate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"username":"InvalidTestUser","email":"","password":"TestPass"}`)
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
	})

	_, err := client.User.Authenticate()
	if err == nil {
		t.Errorf("User.Authenticate with invalid username must return err")
	}
}
