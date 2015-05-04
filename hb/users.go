package hb

import (
	"fmt"
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
	if s.auth == nil {
		return "", fmt.Errorf("credentials are not set")
	}

	const urlStr = "api/v1/users/authenticate"

	req, err := s.client.NewRequest("POST", urlStr, s.auth)
	if err != nil {
		return "", err
	}

	var token string
	_, err = s.client.Do(req, &token)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Get information about a user. Does not require authentication.
func (s *UserService) Get(username string) (*User, error) {
	urlStr := fmt.Sprintf("api/v1/users/%s", url.QueryEscape(username))

	req, err := s.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	user := new(User)
	_, err = s.client.Do(req, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
