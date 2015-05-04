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

// User response looks like this
// {
//   "name": "erengy",
//   "waifu": "Taiga Aisaka",
//   "waifu_or_husbando": "Waifu",
//   "waifu_slug": "toradora",
//   "waifu_char_id": "25930",
//   "location": "",
//   "website": "http://erengy.com",
//   "avatar": "http://static.hummingbird.me/users/avatars/000/002/516/thumb/hb-avatar.jpg?1393289118",
//   "cover_image": "http://static.hummingbird.me/users/cover_images/000/002/516/thumb/hummingbird_cover.jpg?1392287635",
//   "about": null,
//   "bio": "Hi.",
//   "karma": 0,
//   "life_spent_on_anime": 114520,
//   "show_adult_content": true,
//   "title_language_preference": "canonical",
//   "last_library_update": "2014-06-21T19:28:00.443Z",
//   "online": false,
//   "following": false,
//   "favorites": [ *omitted* ]
// }
type User struct {
	Name                    string     `json:"name"`
	Waifu                   string     `json:"waifu"`
	WaifuOrHusbando         string     `json:"waifu_or_husbando"`
	WaifuSlug               string     `json:"waifu_slug"`
	WaifuCharID             int        `json:"waifu_char_id"`
	Location                string     `json:"location"`
	Website                 string     `json:"website"`
	Avatar                  string     `json:"website"`
	CoverImage              string     `json:"cover_image"`
	About                   About      `json:"about"`
	Bio                     string     `json:"bio"`
	Karma                   int        `json:"karma"`
	LifeSpentOnAnime        int        `json:"life_spent_on_anime"`
	ShowAdultContent        bool       `json:"show_adult_content"`
	TitleLanguagePreference string     `json:"title_language_preference"`
	LastLibraryUpdate       time.Time  `json:"last_library_update"`
	Online                  bool       `json:"online"`
	Following               bool       `json:"following"`
	Favorites               []Favorite `json:"favorites"`
}

type About struct {
}

// Favorite response looks like this:
// {
//   "id": 87118,
//   "user_id": 2516,
//   "item_id": 3936,
//   "item_type": "Anime",
//   "created_at": "2014-04-25T11:50:34.831Z",
//   "updated_at": "2014-04-25T11:50:34.831Z",
//   "fav_rank": 9999
// }
type Favorite struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ItemID    int       `json:"item_id"`
	ItemType  string    `json:"item_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FavRank   int       `json:"fav_rank"`
}

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
	const urlStr = "users/authenticate"
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
