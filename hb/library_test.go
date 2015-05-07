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

	libraryEntry, err := client.Library.Update("log-horizon", "valid_user_token", entry)
	if err != nil {
		t.Errorf("Library.Update returned error %v", err)
	}

	got, want := libraryEntry, &LibraryEntry{ID: 7622, EpisodesWatched: 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Library.Update libraryEntry is %v, want %v", got, want)
	}
}

func TestLibraryService_Update_invalidToken(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/libraries/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		requestBody := `{"id":"log-horizon","auth_token":"invalid_user_token"}`
		testBody(t, r, requestBody+"\n")
		testResourceID(t, r, "log-horizon")
		http.Error(w, `{"error": "Invalid authentication token"}`, http.StatusUnauthorized)
	})

	_, err := client.Library.Update("log-horizon", "invalid_user_token", nil)
	if err == nil {
		t.Error("Expected HTTP 401 error.")
	}

	want := fmt.Sprintf("POST %v/api/v1/libraries/log-horizon: 401 Invalid authentication token", server.URL)
	if got := err.Error(); got != want {
		t.Errorf("ErrorResponse is %v, want %v", got, want)
	}
}

func TestLibraryService_Update_badAnimeID(t *testing.T) {
	c := NewClient(nil)
	animeID := "%foo"

	_, err := c.Library.Update(animeID, "", nil)
	if err == nil {
		t.Error("Expected invalid URL escape error.")
	}
}
