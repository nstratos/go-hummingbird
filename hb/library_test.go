package hb

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestLibraryService_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/libraries/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		requestBody := `{"id":"log-horizon","auth_token":"valid_user_token","episodes_watched":3,"increment_episodes":true}`
		testBody(t, r, requestBody+"\n")
		testResourceID(t, r, "log-horizon")
		fmt.Fprintf(w, `{"id":7622,"episodes_watched":4}`)
	})

	entry := &Entry{EpisodesWatched: 3, IncrementEpisodes: true}

	libraryEntry, _, err := client.Library.Update("log-horizon", "valid_user_token", entry)
	if err != nil {
		t.Errorf("Library.Update returned error %v", err)
	}

	got, want := libraryEntry, &LibraryEntry{ID: 7622, EpisodesWatched: 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Library.Update libraryEntry is %v, want %v", got, want)
	}
}

func TestLibraryService_Update_defaultStatus(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/libraries/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		// if no entry is provided, default status "currently-watching" is added.
		requestBody := `{"id":"log-horizon","auth_token":"valid_user_token","status":"currently-watching"}`
		testBody(t, r, requestBody+"\n")
		testResourceID(t, r, "log-horizon")
		fmt.Fprintf(w, `{"id":7622,"status":"currently-watching"}`)
	})

	libraryEntry, _, err := client.Library.Update("log-horizon", "valid_user_token", nil)
	if err != nil {
		t.Errorf("Library.Update returned error %v", err)
	}

	got, want := libraryEntry, &LibraryEntry{ID: 7622, Status: StatusCurrentlyWatching}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Library.Update libraryEntry is %v, want %v", got, want)
	}
}

func TestLibraryService_Update_invalidToken(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/libraries/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		requestBody := `{"id":"log-horizon","auth_token":"invalid_user_token","status":"currently-watching"}`
		testBody(t, r, requestBody+"\n")
		testResourceID(t, r, "log-horizon")
		http.Error(w, `{"error": "Invalid authentication token"}`, http.StatusUnauthorized)
	})

	_, resp, err := client.Library.Update("log-horizon", "invalid_user_token", nil)
	if err == nil {
		t.Error("Expected HTTP 401 error.")
	}

	want := fmt.Sprintf("POST %v/api/v1/libraries/log-horizon: 401 Invalid authentication token", server.URL)
	if got := err.Error(); got != want {
		t.Errorf("ErrorResponse is %v, want %v", got, want)
	}

	if resp == nil {
		t.Error("Expected to return HTTP response despite the API error.")
	}
}

func TestLibraryService_Update_badAnimeID(t *testing.T) {
	c := NewClient(nil)
	animeID := "%foo"

	_, resp, err := c.Library.Update(animeID, "", nil)
	if err == nil {
		t.Error("Expected invalid URL escape error.")
	}

	if resp != nil {
		t.Error("Expected nil HTTP response when NewRequest fails.")
	}
}

func TestLibraryService_Remove(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/libraries/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testResourceParts(t, r, []string{"log-horizon", "remove"})
		requestBody := `{"id":"log-horizon","auth_token":"valid_user_token"}`
		testBody(t, r, requestBody+"\n")
		fmt.Fprintf(w, `true`)
	})

	removed, _, err := client.Library.Remove("log-horizon", "valid_user_token")
	if err != nil {
		t.Errorf("Library.Remove returned error %v", err)
	}

	if got, want := removed, true; got != want {
		t.Errorf("Library.Remove returned %v, want %v", got, want)
	}
}

func TestLibraryService_invalidToken(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/libraries/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testResourceParts(t, r, []string{"log-horizon", "remove"})
		requestBody := `{"id":"log-horizon","auth_token":"invalid_user_token"}`
		testBody(t, r, requestBody+"\n")
		http.Error(w, `{"error": "Invalid authentication token"}`, http.StatusUnauthorized)
	})

	removed, resp, err := client.Library.Remove("log-horizon", "invalid_user_token")
	if err == nil {
		t.Error("Expected HTTP 401 error.")
	}

	if got, want := removed, false; got != want {
		t.Errorf("Library.Remove returned %v, want %v", got, want)
	}

	want := fmt.Sprintf("POST %v/api/v1/libraries/log-horizon/remove: 401 Invalid authentication token", server.URL)
	if got := err.Error(); got != want {
		t.Errorf("ErrorResponse is %v, want %v", got, want)
	}

	if resp == nil {
		t.Error("Expected to return HTTP response despite the API error.")
	}
}

func TestLibraryService_Remove_badAnimeID(t *testing.T) {
	c := NewClient(nil)
	animeID := "%foo"

	_, resp, err := c.Library.Remove(animeID, "")
	if err == nil {
		t.Error("Expected invalid URL escape error.")
	}

	if resp != nil {
		t.Error("Expected nil HTTP response when NewRequest fails.")
	}
}
