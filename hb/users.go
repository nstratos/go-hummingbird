package hb

import (
	"fmt"
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

// UserMini represents a Hummingbird user with minimum info.
// Response looks like:
//   {
//     "name": "erengy",
//     "url": "http://hummingbird.me/users/erengy",
//     "avatar": "http://static.hummingbird.me/users/avatars/000/002/516/thumb/hb-avatar.jpg?1393289118",
//     "avatar_small": "http://static.hummingbird.me/users/avatars/000/002/516/thumb_small/hb-avatar.jpg?1393289118",
//     "nb": false
//   }
type UserMini struct {
	Name        string `json:"name,omitempty"`
	URL         string `json:"url,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	AvatarSmall string `json:"avatar_small,omitempty"`
	NB          bool   `json:"nb,omitempty"`
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
	urlStr := fmt.Sprintf("api/v1/users/%s", username)

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

// Story represents a Hummingbird Story object such as user's activity feed.
// Response looks like:
//   {
//     "id": 2640597,
//     "story_type": "comment",
//     "user": { *omitted* },
//     "updated_at": "2014-06-21T10:37:36.730Z",
//     "self_post": false,
//     "poster": { *omitted* },
//     "substories_count": 1,
//     "substories": [ *omitted* ]
//   }
type Story struct {
	ID              int        `json:"id"`
	StoryType       string     `json:"story_type"`
	User            *UserMini  `json:"user"`
	UpdatedAt       *time.Time `json:"updated_at"`
	SelfPost        bool       `json:"self_post"`
	Poster          *UserMini  `json:"poster"`
	SubstoriesCount int        `json:"substories_count"`
	Substories      []Substory `json:"substories"`
}

// Substory represents a Hummingbird Substory object.
// Response looks like:
//   {
//     "id": 6590163,
//     "substory_type": "watched_episode",
//     "created_at": "2014-06-23T21:25:49.084Z",
//     "episode_number": "12",
//     "service": null,  // should be ignored
//     "permissions": {} // should be ignored
//   }
type Substory struct {
	ID            int        `json:"id"`
	SubstoryType  string     `json:"substory_type"`
	CreatedAt     *time.Time `json:"created_at"`
	EpisodeNumber string     `json:"episode_number"`
}

// Feed returns a user's activity feed.
//
// Does not require authentication.
func (s *UserService) Feed(username string) ([]Story, error) {
	urlStr := fmt.Sprintf("api/v1/users/%s/feed", username)

	req, err := s.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	var stories []Story
	_, err = s.client.Do(req, &stories)
	if err != nil {
		return nil, err
	}
	return stories, nil
}

// FavoriteAnime returns the user's favorite anime in
// an array of Anime objects.
//
// Does not require authentication.
func (s *UserService) FavoriteAnime(username string) ([]Anime, error) {
	urlStr := fmt.Sprintf("api/v1/users/%s/favorite_anime", username)

	req, err := s.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	var anime []Anime
	_, err = s.client.Do(req, &anime)
	if err != nil {
		return nil, err
	}
	return anime, nil
}

// Library returns an array of library entry objects, without genres,
// representing a user's anime library entries.
//
// Does not require authentication.
func (s *UserService) Library(username, status string) ([]LibraryEntry, error) {
	urlStr := fmt.Sprintf("api/v1/users/%s/library", username)

	req, err := s.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	v := req.URL.Query()
	v.Set("status", status)
	req.URL.RawQuery = v.Encode()

	var entries []LibraryEntry
	_, err = s.client.Do(req, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
