package hb

import "time"

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
