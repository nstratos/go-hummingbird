// +build integration

package hb_test

import (
	"flag"
	"fmt"
	"log"
	"testing"

	"github.com/nstratos/go-hummingbird/hb"
)

var (
	hbUsername = flag.String("username", "testgopher", "Hummingbird.me username to use for integration tests")
	hbPassword = flag.String("password", "", "Hummingbird.me password to use for integration tests")

	client *hb.Client
)

func TestAnimeServiceIntegration(t *testing.T) {

	c := hb.NewClient(nil)

	// Acquire user's authentication token.
	//token, _, err := c.User.Authenticate(*hbUsername, "", *hbPassword)
	//if err != nil {
	//	log.Fatalf("token err = %v", err)
	//}

	anime, _, err := c.Anime.Get("nichijou", "english")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Title:", anime.Title)
	fmt.Println("Synopsis:", anime.Synopsis)
}
