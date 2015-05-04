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

// Client manages communication with the Hummingbird API.
type Client struct {
	client *http.Client

	BaseURL *url.URL

	User *UserService
}

// NewClient returns a new Hummingbird API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL}

	c.User = &UserService{client: c}
	return c
}
