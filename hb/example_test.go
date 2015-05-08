package hb_test

import (
	"fmt"
	"log"

	"bitbucket.org/nstratos/go-hummingbird/hb"
)

func ExampleUserService_Authenticate() {
	c := hb.NewClient(nil)
	token, _, err := c.User.Authenticate("USER_USERNAME", "", "USER_PASSWORD")
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
