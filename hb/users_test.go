package hb

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUserService_Authenticate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/users/authenticate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"username":"TestUser","email":"","password":"TestPass"}`+"\n")
		fmt.Fprintf(w, `"token1234"`)
	})

	token, _, err := client.User.Authenticate("TestUser", "", "TestPass")
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

	mux.HandleFunc("/api/v1/users/authenticate", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBody(t, r, `{"username":"InvalidTestUser","email":"","password":"TestPass"}`+"\n")
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
	})

	_, resp, err := client.User.Authenticate("InvalidTestUser", "", "TestPass")
	if err == nil {
		t.Error("Expected HTTP 401 error.")
	}

	if resp == nil {
		t.Error("Expected to return HTTP response despite the API error.")
	}
}

func TestUserService_Authenticate_credentialsNotSet(t *testing.T) {
	c := NewClient(nil)
	_, resp, err := c.User.Authenticate("", "", "")
	if err == nil {
		t.Errorf("Expected username or email must be provided error.")
	}

	if resp != nil {
		t.Error("Expected nil HTTP response when credentials not set error.")
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

	user, _, err := client.User.Get("TestUser")
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

	_, resp, err := client.User.Get("TestUser")
	if err == nil {
		t.Error("Expected HTTP 404 error.")
	}

	if resp == nil {
		t.Error("Expected to return HTTP response despite the API error.")
	}
}

func TestUserService_Get_badUsername(t *testing.T) {
	c := NewClient(nil)
	username := "%foo"

	_, resp, err := c.User.Get(username)
	if err == nil {
		t.Error("Expected invalid URL escape error.")
	}

	if resp != nil {
		t.Error("Expected nil HTTP response when NewRequest fails.")
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

	stories, _, err := client.User.Feed("TestUser")
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

	_, resp, err := client.User.Feed("TestUser")
	if err == nil {
		t.Error("Expected HTTP 404 error.")
	}

	if resp == nil {
		t.Error("Expected to return HTTP response despite the API error.")
	}
}

func TestUserService_Feed_badUsername(t *testing.T) {
	c := NewClient(nil)
	username := "%foo"

	_, resp, err := c.User.Feed(username)
	if err == nil {
		t.Error("Expected invalid URL escape error.")
	}

	if resp != nil {
		t.Error("Expected nil HTTP response when NewRequest fails.")
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

	anime, _, err := client.User.FavoriteAnime("TestUser")
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

	_, resp, err := client.User.FavoriteAnime("InvalidUser")
	if err == nil {
		t.Error("Expected HTTP 404 error.")
	}

	if resp == nil {
		t.Error("Expected to return HTTP response despite the API error.")
	}
}

func TestUserService_FavoriteAnime_badUsername(t *testing.T) {
	c := NewClient(nil)
	username := "%foo"

	_, resp, err := c.User.FavoriteAnime(username)
	if err == nil {
		t.Error("Expected invalid URL escape error.")
	}

	if resp != nil {
		t.Error("Expected nil HTTP response when NewRequest fails.")
	}
}

func TestUserService_Library(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/users/TestUser/library", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"status": "currently-watching"})
		testResourceParts(t, r, []string{"TestUser", "library"})
		fmt.Fprintf(w, `
			[
			  { "id":22, "status":"currently-watching", 
			    "rating":{"type":"advanced", "value":"4.0"} 
			  },
			  { "id":23, "status":"currently-watching", 
			    "anime":{"title":"Log Horizon"} 
			  }
			]
			`)
	})

	entries, _, err := client.User.Library("TestUser", "currently-watching")
	if err != nil {
		t.Errorf("User.Library returned error %v", err)
	}

	got, want := entries, []LibraryEntry{
		{ID: 22, Status: StatusCurrentlyWatching,
			Rating: &LibraryEntryRating{Type: "advanced", Value: "4.0"},
		},
		{ID: 23, Status: StatusCurrentlyWatching,
			Anime: &Anime{Title: "Log Horizon"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("User.Library entries are %v, want %v", got, want)
	}
}

func TestUserService_Library_notFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/users/InvalidTestUser/library", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"status": "currently-watching"})
		testResourceParts(t, r, []string{"InvalidTestUser", "library"})
		http.Error(w, "not found", http.StatusNotFound)
	})

	_, resp, err := client.User.Library("InvalidTestUser", "currently-watching")
	if err == nil {
		t.Error("Expected HTTP 404 error.")
	}

	if resp == nil {
		t.Error("Expected to return HTTP response despite the API error.")
	}
}

func TestUserService_Library_badUsername(t *testing.T) {
	c := NewClient(nil)
	username := "%foo"

	_, resp, err := c.User.Library(username, "")
	if err == nil {
		t.Error("Expected invalid URL escape error.")
	}

	if resp != nil {
		t.Error("Expected nil HTTP response when NewRequest fails.")
	}
}
