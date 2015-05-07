package hb

import (
	"fmt"
	"time"
)

// LibraryEntry represents a library entry of a Hummingbird user.
// Response looks like:
//   {
//     "id": 5593549,
//     "episodes_watched": 20,
//     "last_watched": "2014-06-20T05:31:27.074Z",
//     "updated_at": "2014-08-18T16:04:05.383Z",
//     "rewatched_times": 0,
//     "notes": "",
//     "notes_present": false,
//     "status": "currently-watching",
//     "private": false,
//     "rewatching": false,
//     "anime": { *omitted* },
//     "rating": { *omitted* }
//   }
type LibraryEntry struct {
	ID              int                 `json:"id"`
	EpisodesWatched int                 `json:"episodes_watched"`
	LastWatched     *time.Time          `json:"last_watched"`
	UpdatedAt       *time.Time          `json:"updated_at"`
	RewatchedTimes  int                 `json:"rewatched_times"`
	Notes           string              `json:"notes"`
	NotesPresent    bool                `json:"notes_present"`
	Status          string              `json:"status"`
	Private         bool                `json:"private"`
	Rewatching      bool                `json:"rewatching"`
	Anime           *Anime              `json:"anime"`
	Rating          *LibraryEntryRating `json:"rating"`
}

// LibraryEntryRating represents the rating of a user's library entry.
// The representation it's value depends on the type which can be either
// "simple" or "advanced".
//
// If type is "simple", the value can be "negative", "neutral" or "positive".
// If type is "advanced", the value can be a number between "0.0" and "5.0".
//
// For conversion between "simple" and "advanced":
//   0   <= "negative" => 2.4
//   2.4 <  "neutral"   > 3.6
//   3.6 <= "positive" => 5
type LibraryEntryRating struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// LibraryServices handles communication with the Hummingbird API library
// methods (GET /users/{username}/library} is handled by UserService).
//
// Hummingbird API docs:
// https://github.com/hummingbird-me/hummingbird/wiki/API-v1-Methods#library
type LibraryService struct {
	client *Client
}

// Entry represents the values that the user can add/update on a Library entry.
//
// ID - Required
//
// Can be an anime ID like "7622" or a slug like "log-horizon".
//
// AuthToken - Required
//
// A valid user authentication token. One way to acquire a token is:
//   c := hb.NewClient(nil)
//   c.SetCredentials("YOUR_HUMMINGBIRD_USERNAME", "", "YOUR_HUMMINGBIRD_PASSWORD")
//   token, err := c.User.Authenticate()
//   // handle err
// Note that c.User.Authenticate(), if successful, keeps the token in memory.
// That token will be used automatically, if the Entry.AuthToken value is empty.
//
// Status - Optional
//
// Can be one of:
//   hb.StatusCurrentlyWatching
//   hb.StatusPlanToWatch
//   hb.StatusCompleted
//   hb.StatusOnHold
//   hb.StatusDropped
//
// Privacy - Optional
//
// Can be either "public" or "private". Making an entry private will hide it
// from public view.
//
// Rating - Optional
//
// Can be one of:
//   "0", "0.5", "1", "1.5", "2", "2.5", "3", "3.5", "4", "4.5", "5".
// Setting it to the current value or "0" will remove the rating.
//
// SaneRatingUpdate - Optional
//
// Can be one of:
//   "0", "0.5", "1", "1.5", "2", "2.5", "3", "3.5", "4", "4.5", "5".
// Setting it to "0" will remove the rating. This should be used instead of
// Rating if you don't want to unset the rating when setting it to the value
// it already has.
//
// Rewatching - Optional
//
// Can be true or false.
//
// RewatchedTimes - Optional
//
// Number of rewatches. Can be 0 or above.
//
// Notes - Optional
//
// Personal notes.
//
// EpisodesWatched - Optional
//
// Number of watched episodes. Can be between 0 and the total number of episodes.
// If equal to total number of episodes, Status should be set to "completed".
//
// IncrementEpisodes - Optional
//
// If set to true, increments the number of watched episodes by one.
// If used along with EpisodesWatched, provided value will be incremented.
//
// Hummingbird API docs:
// https://github.com/hummingbird-me/hummingbird/wiki/API-v1-Methods#parameters-3
type Entry struct {
	ID                string `json:"id"`
	AuthToken         string `json:"auth_token"`
	Status            string `json:"status,omitempty"`
	Privacy           string `json:"privacy,omitempty"`
	Rating            string `json:"rating,omitempty"`
	SaneRatingUpdate  string `json:"sane_rating_update,omitempty"`
	Rewatching        bool   `json:"rewatching,omitempty"`
	RewatchedTimes    int    `json:"rewatched_times,omitempty"`
	Notes             string `json:"notes,omitempty"`
	EpisodesWatched   int    `json:"episodes_watched,omitempty"`
	IncrementEpisodes bool   `json:"increment_episodes,omitempty"`
}

// Update adds or updates a user's library entry. The updated Library entry is
// returned on success.
//
// It assumes that the client already has a valid auth token by a previous
// successful authentication. Alternatively an auth token can be directly
// passed in the entry values.
//
// Requires authentication.
func (s *LibraryService) Update(entry Entry) (*LibraryEntry, error) {
	urlStr := fmt.Sprintf("api/v1/libraries/%v", entry.ID)

	if entry.AuthToken == "" {
		if auth := s.client.User.auth; auth != nil && auth.token != "" {
			entry.AuthToken = auth.token
		}
	}

	req, err := s.client.NewRequest("POST", urlStr, entry)
	if err != nil {
		return nil, err
	}

	libraryEntry := new(LibraryEntry)
	_, err = s.client.Do(req, libraryEntry)
	if err != nil {
		return nil, err
	}
	return libraryEntry, nil
}

// func (c *Client) fetchAuthToken() (string, error) {

// 	if c.User.auth != nil && c.User.auth.token != "" {
// 		return c.User.auth.token, nil
// 	}

// 	token, err := c.User.Authenticate()
// 	if err != nil {
// 		return "", err
// 	}
// 	c.User.auth.token = token

// 	return token, nil
// }
