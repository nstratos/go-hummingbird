package hb

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
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

	mux.HandleFunc("/api/v1/users/authenticate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"username":"TestUser","email":"","password":"TestPass"}`+"\n")
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

	mux.HandleFunc("/api/v1/users/authenticate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"username":"InvalidTestUser","email":"","password":"TestPass"}`+"\n")
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
	})

	_, err := client.User.Authenticate()
	if err == nil {
		t.Errorf("User.Authenticate with invalid username must return err")
	}
}

func TestUserService_Authenticate_credentialsNotSet(t *testing.T) {
	c := NewClient(nil)
	_, err := c.User.Authenticate()
	if err == nil {
		t.Errorf("Expected credentials not set error.")
	}
}

func TestUserService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/users/TestUser", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testResourceID(t, r, "TestUser")
		fmt.Fprintf(w, `{"name":"TestUser","bio":"My bio."}`)
	})

	user, err := client.User.Get("TestUser")
	if err != nil {
		t.Errorf("User.Get returned error %v", err)
	}

	got, want := user, &User{Name: "TestUser", Bio: "My bio."}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("User.Get user is %v, want %v", got, want)
	}
}

func TestUserService_Get_notFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/users/TestUser", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testResourceID(t, r, "TestUser")
		http.Error(w, "not found", http.StatusNotFound)
	})

	_, err := client.User.Get("TestUser")
	if err == nil {
		t.Error("Expected HTTP 404 error.")
	}
}

func TestUserService_Get_badURL(t *testing.T) {
	c := NewClient(nil)
	urlStr := "%foo"

	_, err := c.User.Get(urlStr)
	if err == nil {
		t.Error("Expected invalid URL escape error.")
	}
}

func TestUserService_Feed(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/users/TestUser/feed", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testResourceParts(t, r, []string{"TestUser", "feed"})
		fmt.Fprintf(w, `[{"id":1,"story_type":"comment"},{"id":2,"story_type":"media_story"}]`)
	})

	stories, err := client.User.Feed("TestUser")
	if err != nil {
		t.Errorf("User.Feed returned error %v", err)
	}

	got, want := stories, []Story{{ID: 1, StoryType: "comment"}, {ID: 2, StoryType: "media_story"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("User.Feed stories are %v, want %v", got, want)
	}
}

func TestUserService_Feed_notFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/users/InvalidUser/feed", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testResourceParts(t, r, []string{"InvalidUser", "feed"})
		http.Error(w, "not found", http.StatusNotFound)
	})

	_, err := client.User.Feed("TestUser")
	if err == nil {
		t.Error("Expected HTTP 404 error.")
	}
}

func TestUserService_Feed_badURL(t *testing.T) {
	c := NewClient(nil)
	urlStr := "%foo"

	_, err := c.User.Feed(urlStr)
	if err == nil {
		t.Error("Expected invalid URL escape error.")
	}
}

func TestUserService_FavoriteAnime(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/users/TestUser/favorite_anime", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testResourceParts(t, r, []string{"TestUser", "favorite_anime"})
		fmt.Fprintf(w, `[{"title":"Log Horizon"},{"title":"Nichijou"}]`)
	})

	anime, err := client.User.FavoriteAnime("TestUser")
	if err != nil {
		t.Errorf("User.FavoriteAnime returned error %v", err)
	}

	got, want := anime, []Anime{{Title: "Log Horizon"}, {Title: "Nichijou"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("User.FavoriteAnime anime are %v, want %v", got, want)
	}
}

func TestUserService_FavoriteAnime_notFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/users/InvalidUser/favorite_anime", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testResourceParts(t, r, []string{"InvalidUser", "favorite_anime"})
		http.Error(w, "not found", http.StatusNotFound)
	})

	_, err := client.User.FavoriteAnime("InvalidUser")
	if err == nil {
		t.Error("Expected HTTP 404 error.")
	}
}

func TestUserService_FavoriteAnime_badURL(t *testing.T) {
	c := NewClient(nil)
	urlStr := "%foo"

	_, err := c.User.FavoriteAnime(urlStr)
	if err == nil {
		t.Error("Expected invalid URL escape error.")
	}
}
