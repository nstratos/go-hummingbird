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
// Can be an anime ID like "7622" or a slug like "log-horizon". It is set
// automatically from the method that uses Entry.
//
// AuthToken - Required
//
// A valid user authentication token. It is set automatically from the method
// that uses Entry.
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

// Update adds or updates a user's library entry. The updated library entry is
// returned on success. Requires authentication.
//
// The animeID can be an ID like "7622" or a slug like "log-horizon".
//
// To acquire a user's authentication token:
//   c := hb.NewClient(nil)
//   token, err := c.User.Authenticate("USER_HUMMINGBIRD_USERNAME", "", "USER_HUMMINGBIRD_PASSWORD")
//   // handle err
//
// An optional entry parameter can be specified with additional values to
// add/update on a user's library entry.
func (s *LibraryService) Update(animeID, authToken string, entry *Entry) (*LibraryEntry, error) {
	urlStr := fmt.Sprintf("api/v1/libraries/%v", animeID)

	if entry == nil {
		entry = new(Entry)
	}
	entry.ID = animeID
	entry.AuthToken = authToken

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
