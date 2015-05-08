package hb

import (
	"fmt"
	"net/http"
	"time"
)

// User represents a Hummingbird user.
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
type UserMini struct {
	Name        string `json:"name,omitempty"`
	URL         string `json:"url,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	AvatarSmall string `json:"avatar_small,omitempty"`
	NB          bool   `json:"nb,omitempty"`
}

// Favorite represents a favorite item of a Hummingbird user.
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
}

type auth struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Authenticate a user and return an authentication token if successful. That
// token can be used in other methods that require authentication. From
// username and email only one is needed.
func (s *UserService) Authenticate(username, email, password string) (string, *http.Response, error) {
	if username == "" && email == "" {
		return "", nil, fmt.Errorf("username or email must be provided")
	}

	const urlStr = "api/v1/users/authenticate"

	auth := auth{Username: username, Email: email, Password: password}

	req, err := s.client.NewRequest("POST", urlStr, auth)
	if err != nil {
		return "", nil, err
	}

	var token string
	resp, err := s.client.Do(req, &token)
	if err != nil {
		return "", resp, err
	}

	return token, resp, nil
}

// Get information about a user.
//
// Does not require authentication.
func (s *UserService) Get(username string) (*User, *http.Response, error) {
	urlStr := fmt.Sprintf("api/v1/users/%s", username)

	req, err := s.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}
	return user, resp, nil
}

// Story represents a Hummingbird Story object such as user's activity feed.
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
type Substory struct {
	ID            int        `json:"id"`
	SubstoryType  string     `json:"substory_type"`
	CreatedAt     *time.Time `json:"created_at"`
	EpisodeNumber string     `json:"episode_number"`
}

// Feed returns a user's activity feed.
//
// Does not require authentication.
func (s *UserService) Feed(username string) ([]Story, *http.Response, error) {
	urlStr := fmt.Sprintf("api/v1/users/%s/feed", username)

	req, err := s.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, nil, err
	}

	var stories []Story
	resp, err := s.client.Do(req, &stories)
	if err != nil {
		return nil, resp, err
	}
	return stories, resp, nil
}

// FavoriteAnime returns the user's favorite anime in
// an array of Anime objects.
//
// Does not require authentication.
func (s *UserService) FavoriteAnime(username string) ([]Anime, *http.Response, error) {
	urlStr := fmt.Sprintf("api/v1/users/%s/favorite_anime", username)

	req, err := s.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, nil, err
	}

	var anime []Anime
	resp, err := s.client.Do(req, &anime)
	if err != nil {
		return nil, resp, err
	}
	return anime, resp, nil
}

// Library returns an array of library entry objects, without genres,
// representing a user's anime library entries.
//
// Does not require authentication.
func (s *UserService) Library(username, status string) ([]LibraryEntry, *http.Response, error) {
	urlStr := fmt.Sprintf("api/v1/users/%s/library", username)

	req, err := s.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, nil, err
	}

	v := req.URL.Query()
	v.Set("status", status)
	req.URL.RawQuery = v.Encode()

	var entries []LibraryEntry
	resp, err := s.client.Do(req, &entries)
	if err != nil {
		return nil, resp, err
	}
	return entries, resp, nil
}
