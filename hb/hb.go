package hb

import (
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "http://hummingbird.me/api/v1/"
)

type auth struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Client struct {
	client *http.Client

	BaseURL *url.URL

	User *UserService
}

func NewClient(httpClient *http.Client) *Client {
	if client == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL}

	c.User = &UserService{client: c}
	return c
}
