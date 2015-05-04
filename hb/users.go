package hb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// User represents a Hummingbird user.
// Response looks like:
//   {
//     "name": "erengy",
//     "waifu": "Taiga Aisaka",
//     "waifu_or_husbando": "Waifu",
//     "waifu_slug": "toradora",
//     "waifu_char_id": "25930",
//     "location": "",
//     "website": "http://erengy.com",
//     "avatar": "http://static.hummingbird.me/users/avatars/000/002/516/thumb/hb-avatar.jpg?1393289118",
//     "cover_image": "http://static.hummingbird.me/users/cover_images/000/002/516/thumb/hummingbird_cover.jpg?1392287635",
//     "about": null,
//     "bio": "Hi.",
//     "karma": 0,
//     "life_spent_on_anime": 114520,
//     "show_adult_content": true,
//     "title_language_preference": "canonical",
//     "last_library_update": "2014-06-21T19:28:00.443Z",
//     "online": false,
//     "following": false,
//     "favorites": [ *omitted* ]
//   }
type User struct {
	Name                    string     `json:"name,omitempty"`
	Waifu                   string     `json:"waifu,omitempty"`
	WaifuOrHusbando         string     `json:"waifu_or_husbando,omitempty"`
	WaifuSlug               string     `json:"waifu_slug,omitempty"`
	WaifuCharID             string     `json:"waifu_char_id,omitempty"`
	Location                string     `json:"location,omitempty"`
	Website                 string     `json:"website,omitempty"`
	Avatar                  string     `json:"website,omitempty"`
	CoverImage              string     `json:"cover_image,omitempty"`
	About                   string     `json:"about,omitempty"`
	Bio                     string     `json:"bio,omitempty"`
	Karma                   int        `json:"karma,omitempty"`
	LifeSpentOnAnime        int        `json:"life_spent_on_anime,omitempty"`
	ShowAdultContent        bool       `json:"show_adult_content,omitempty"`
	TitleLanguagePreference string     `json:"title_language_preference,omitempty"`
	LastLibraryUpdate       *time.Time `json:"last_library_update,omitempty"`
	Online                  bool       `json:"online,omitempty"`
	Following               bool       `json:"following,omitempty"`
	Favorites               []Favorite `json:"favorites,omitempty"`
}

// Favorite represents a favorite item of a Hummingbird user.
// Response looks like:
//   {
//     "id": 87118,
//     "user_id": 2516,
//     "item_id": 3936,
//     "item_type": "Anime",
//     "created_at": "2014-04-25T11:50:34.831Z",
//     "updated_at": "2014-04-25T11:50:34.831Z",
//     "fav_rank": 9999
//   }
type Favorite struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ItemID    int       `json:"item_id"`
	ItemType  string    `json:"item_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FavRank   int       `json:"fav_rank"`
}

// UserService handles communication with the user methods of
// the Hummingbird API.
//
// Hummingbird API docs:
// https://github.com/hummingbird-me/hummingbird/wiki/API-v1-Methods#user
type UserService struct {
	client *Client
	auth   *auth
}

// SetCredentials sets the username, email and password for the methods that
// require authentication. From username and email only one is needed.
func (s *UserService) SetCredentials(username, email, password string) {
	s.auth = &auth{username, email, password}
}

// Authenticate a user and return an authentication token if successful. That
// authentication token can be used in other methods that require authentication.
func (s *UserService) Authenticate() (string, error) {
	const urlStr = "api/v1/users/authenticate"
	endpoint, _ := url.Parse(urlStr)
	u := s.client.BaseURL.ResolveReference(endpoint)

	if s.auth == nil {
		return "", fmt.Errorf("credentials are not set")
	}
	b, err := json.Marshal(s.auth)
	if err != nil {
		return "", fmt.Errorf("cannot marshal auth: %v", err)
	}
	resp, err := http.Post(u.String(), "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read body")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("unauthorized")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("%v %v: %v", resp.Request.Method, resp.Request.URL, resp.StatusCode)
	}
	var token string
	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal token: %v", err)
	}
	return token, nil
}

// Get information about a user. Does not require authentication.
func (s *UserService) Get(username string) (*User, error) {
	endpoint, _ := url.Parse(fmt.Sprintf("api/v1/users/%s", username))
	u := s.client.BaseURL.ResolveReference(endpoint)

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body")
	}
	defer resp.Body.Close()

	user := new(User)
	err = json.Unmarshal(body, user)
	if err != nil {
		return nil, fmt.Errorf("received: %v, cannot unmarshal user: %v", string(body), err)
	}
	return user, nil
}
