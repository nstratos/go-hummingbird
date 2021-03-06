package hb_test

import (
	"fmt"
	"log"

	"github.com/nstratos/go-hummingbird/hb"
)

func ExampleAnimeService_Get() {
	c := hb.NewClient(nil)

	anime, _, err := c.Anime.Get("nichijou", "english")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Title:", anime.Title)
	fmt.Println("Synopsis:", anime.Synopsis)
}

func ExampleAnimeService_Search() {
	c := hb.NewClient(nil)

	anime, _, err := c.Anime.Search("anohana")
	if err != nil {
		log.Fatal(err)
	}
	for i, a := range anime {
		fmt.Printf("--- Search result #%d ---\n", i+1)
		fmt.Printf("Title: %v\n", a.Title)
		fmt.Printf("Synopsis: %v\n\n", a.Synopsis)
	}
}

func ExampleUserService_Authenticate() {
	c := hb.NewClient(nil)

	token, _, err := c.User.Authenticate("", "USER_HUMMINGBIRD_EMAIL", "USER_HUMMINGBIRD_PASSWORD")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user auth token:", token)
}

func ExampleUserService_Get() {
	c := hb.NewClient(nil)

	u, _, err := c.User.Get("cybrox")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Name:", u.Name)
	fmt.Println("About:", u.About)
}

func ExampleUserService_FavoriteAnime() {
	c := hb.NewClient(nil)

	anime, _, err := c.User.FavoriteAnime("cybrox")
	if err != nil {
		log.Fatal(err)
	}
	for i, a := range anime {
		fmt.Printf("--- Favorite anime #%d ---\n", i+1)
		fmt.Println("Title:", a.Title)
		fmt.Println("Synopsis:", a.Synopsis)
	}
}

func ExampleUserService_Library() {
	c := hb.NewClient(nil)

	entries, _, err := c.User.Library("cybrox", hb.StatusCurrentlyWatching)
	if err != nil {
		log.Fatal(err)
	}
	for i, e := range entries {
		fmt.Printf("--- Library entry #%d ---\n", i+1)
		fmt.Printf("Anime Title: %v\n", e.Anime.Title)
		fmt.Printf("Episodes watched: %v\n\n", e.EpisodesWatched)
	}
}

func ExampleUserService_Library_getAllStatuses() {
	c := hb.NewClient(nil)

	entries, _, err := c.User.Library("cybrox", "")
	if err != nil {
		log.Fatal(err)
	}
	for i, e := range entries {
		fmt.Printf("--- Library entry #%d ---\n", i+1)
		fmt.Printf("Anime Title: %v\n", e.Anime.Title)
		fmt.Printf("Episodes watched: %v\n\n", e.EpisodesWatched)
	}
}

func ExampleUserService_Feed() {
	c := hb.NewClient(nil)

	stories, _, err := c.User.Feed("cybrox")
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range stories {
		if s.StoryType == "comment" {
			for _, ss := range s.Substories {
				if ss.SubstoryType == "comment" {
					fmt.Printf("%v said to %v:\n", s.Poster.Name, s.User.Name)
					fmt.Printf("%v\n\n", ss.Comment)
				}
			}
		}
	}
}

func ExampleLibraryService_Update() {
	c := hb.NewClient(nil)

	// Acquire user's authentication token.
	token, _, err := c.User.Authenticate("USER_HUMMINGBIRD_USERNAME", "", "USER_HUMMINGBIRD_PASSWORD")
	if err != nil {
		log.Fatalf("token err = %v", err)
	}

	// Quick function to keep example shorter. For real code, always remember:
	// "Don't just check errors, handle them gracefully."
	checkErr := func(err error) {
		if err != nil {
			log.Println("Update failed:", err)
		}
	}

	// Add anime Nichijou to the user's library. Providing nil will still add
	// Status "Currently watching" by default.
	_, _, err = c.Library.Update("nichijou", token, nil)
	checkErr(err)

	// Update Nichijou, increasing episodes watched by one.
	_, _, err = c.Library.Update("nichijou", token, &hb.Entry{IncrementEpisodes: true})
	checkErr(err)

	// Update Nichijou, setting status as on-hold.
	_, _, err = c.Library.Update("nichijou", token, &hb.Entry{Status: hb.StatusOnHold})
	checkErr(err)

	// Update Nichijou, setting status as currently-watching and number of
	// episodes watched as 5.
	_, _, err = c.Library.Update("nichijou", token, &hb.Entry{
		Status:          hb.StatusCurrentlyWatching,
		EpisodesWatched: 5,
	})
	checkErr(err)

	// Update Nichijou, setting status as completed and setting a note.
	_, _, err = c.Library.Update("nichijou", token, &hb.Entry{
		Status: hb.StatusCompleted,
		Notes:  "crazy",
	})
	checkErr(err)

}

func ExampleLibraryService_Remove() {
	c := hb.NewClient(nil)

	// Acquire user's authentication token.
	token, _, err := c.User.Authenticate("USER_HUMMINGBIRD_USERNAME", "", "USER_HUMMINGBIRD_PASSWORD")
	if err != nil {
		log.Fatalf("token err = %v", err)
	}

	removed, _, err := c.Library.Remove("nichijou", token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Anime was removed:", removed)
}
