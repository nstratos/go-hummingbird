package hb

import (
	"fmt"
	"net/url"
)

// Anime represents a hummingbird anime object.
// An anime response looks like:
//   {
//     "id": 7622,
//     "slug": "log-horizon",
//     "status": "Finished Airing",
//     "url": "https://hummingbird.me/anime/log-horizon",
//     "title": "Log Horizon",
//     "alternate_title": "",
//     "episode_count": 25,
//     "episode_length": 25,
//     "cover_image": "https://static.hummingbird.me/anime/poster_images/000/007/622/large/b0012149_5229cf3c7f4ee.jpg?1408461927",
//     "synopsis": "The story begins when 30,000 Japanese gamers are trapped in the fantasy online game world Elder Tale. What was once a sword-and-sorcery world is now the real world. The main lead Shiroe attempts to survive with his old friend Naotsugu and the beautiful assassin Akatsuki.\r\n(Source: ANN)",
//     "show_type": "TV",
//     "started_airing": "2013-10-05",
//     "finished_airing": "2014-03-22",
//     "community_rating": 4.16741419054807,
//     "age_rating": "PG13",
//     "genres": [
//       { "name": "Action" },
//       { "name": "Adventure" },
//       { "name": "Magic" },
//       { "name": "Fantasy" },
//       { "name": "Game" }
//     ]
//   }
type Anime struct {
	ID              int     `json:"id,omitempty"`
	Slug            string  `json:"slug,omitempty"`
	Status          string  `json:"status,omitempty"`
	URL             string  `json:"url,omitempty"`
	Title           string  `json:"title,omitempty"`
	AlternateTitle  string  `json:"alternate_title,omitempty"`
	EpisodeCount    int     `json:"episode_count,omitempty"`
	EpisodeLength   int     `json:"episode_length,omitempty"`
	CoverImage      string  `json:"cover_image,omitempty"`
	Synopsis        string  `json:"synopsis,omitempty"`
	ShowType        string  `json:"show_type,omitempty"`
	StartedAiring   string  `json:"started_airing,omitempty"`
	FinishedAiring  string  `json:"finished_airing,omitempty"`
	CommunityRating float64 `json:"community_rating,omitempty"`
	AgeRating       string  `json:"age_rating,omitempty"`
	Genres          []Genre `json:"genres,omitempty"`
	FavID           int     `json:"fav_id,omitempty"`   // When requesting user favorite anime.
	FavRank         int     `json:"fav_rank,omitempty"` // When requesting user favorite anime.
}

// Genre represents the genre of an anime.
type Genre struct {
	Name string
}

// AnimeService handles communication with the anime methods of
// the Hummingbird API.
//
// Hummingbird API docs:
// https://github.com/hummingbird-me/hummingbird/wiki/API-v1-Methods#anime
type AnimeService struct {
	client *Client
}

// Get returns anime metadata based on ID which can be either the anime ID or
// a slug. An optional parameter about the title language preference can be
// used which can be one of: "canonical", "english", "romanized".
// If omitted, "canonical" will be used.
//
// Does not require authentication.
func (s *AnimeService) Get(animeID, titleLangPref string) (*Anime, error) {
	urlStr := fmt.Sprintf("api/v1/anime/%s", url.QueryEscape(animeID))

	req, err := s.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	if titleLangPref != "" {
		v := req.URL.Query()
		v.Set("title_language_preference", titleLangPref)
		req.URL.RawQuery = v.Encode()
	}

	anime := new(Anime)
	_, err = s.client.Do(req, anime)
	if err != nil {
		return nil, err
	}
	return anime, nil
}

// Search allows searching anime by title. It returns an array of anime objects
// (5 max) without genres. It supports fuzzy search.
//
// Does not require authentication.
func (s *AnimeService) Search(query string) ([]Anime, error) {
	const urlStr = "api/v1/search/anime"

	req, err := s.client.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	v := req.URL.Query()
	v.Set("query", query)
	req.URL.RawQuery = v.Encode()

	var anime []Anime
	_, err = s.client.Do(req, &anime)
	if err != nil {
		return nil, err
	}
	return anime, nil
}
