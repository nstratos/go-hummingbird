package hb

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAnimeService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/anime/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testResourceID(t, r, "log-horizon")
		testFormValues(t, r, values{"title_language_preference": "english"})
		fmt.Fprintf(w, `{"title":"Log Horizon"}`)
	})

	anime, err := client.Anime.Get("log-horizon", "english")
	if err != nil {
		t.Errorf("Anime.Get returned error %v", err)
	}

	got, want := anime, &Anime{Title: "Log Horizon"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.Get anime is %v, want %v", got, want)
	}
}

func TestAnimeService_Get_notFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/anime/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testResourceID(t, r, "invalid-anime")
		testFormValues(t, r, values{})
		http.Error(w, "not found", http.StatusNotFound)
	})

	_, err := client.Anime.Get("invalid-anime", "")
	if err == nil {
		t.Errorf("Expected HTTP 404 error.")
	}
}

func TestAnimeService_Search(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/search/anime", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"query": "log horizon"})
		fmt.Fprintf(w, `[{"title":"Log Horizon1"},{"title":"Log Horizon2"}]`)
	})

	result, err := client.Anime.Search("log horizon")
	if err != nil {
		t.Errorf("Anime.Search returned error %v", err)
	}

	got, want := result, []Anime{{Title: "Log Horizon1"}, {Title: "Log Horizon2"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Anime.Search result is %v, want %v", got, want)
	}
}

func TestAnimeService_Search_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/search/anime", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"query": "log horizon"})
		http.Error(w, "something broke", http.StatusInternalServerError)
	})

	_, err := client.Anime.Search("log horizon")
	if err == nil {
		t.Errorf("Expected HTTP 500 error.")
	}
}
