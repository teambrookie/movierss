package trakt

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/teambrookie/movierss/dao"
)

type Movie dao.Movie

type MovieProvider interface {
	Collection(slug string) ([]Movie, error)
	WatchList(slug string, filter string) ([]Movie, error)
}

type Trakt struct {
	APIKey string
}

type trackResponse []struct {
	Rank  int    `json:"rank"`
	Type  string `json:"type"`
	Movie Movie
}

func respToMovies(r trackResponse) []Movie {
	var movies []Movie
	for _, m := range r {
		movies = append(movies, m.Movie)
	}
	return movies
}

func diff(watchlist, collection []Movie) []Movie {
	var rest []Movie
	isHere := false
	for _, w := range watchlist {
		for _, c := range collection {
			if w.Ids.Imdb == c.Ids.Imdb {
				isHere = true
				break
			}
		}
		if !isHere {
			rest = append(rest, w)
		}
		isHere = false
	}
	return rest
}

// Collection return the content of your trakt.tv collection
func (p Trakt) Collection(slug string) ([]Movie, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.trakt.tv/users/%s/collection/movies", slug)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("trakt-api-version", "2")
	req.Header.Add("trakt-api-key", "53d18e80e9c0104d27253918a69a4a9c9e8d47e2c4dbf6aaf98d34d2f1f558a4")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return nil, err
	}

	var response trackResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("Error when decoding response")
		return nil, err
	}
	return respToMovies(response), nil
}

// WatchList return the content of your trakt.tv watchlist
// filter can be set to notCollected or ""
func (p Trakt) WatchList(slug string, filter string) ([]Movie, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.trakt.tv/users/%s/watchlist/movie", slug)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("trakt-api-version", "2")
	req.Header.Add("trakt-api-key", "53d18e80e9c0104d27253918a69a4a9c9e8d47e2c4dbf6aaf98d34d2f1f558a4")

	resp, err := client.Do(req)

	if err != nil {
		log.Println("Errored when sending request to the server")
		return nil, err
	}
	defer resp.Body.Close()

	var response trackResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("Error when decoding response")
		return nil, err
	}

	watchlist := respToMovies(response)

	if filter == "notCollected" {
		collection, err := Collection(slug)
		if err != nil {
			log.Println("Error when querrying for trakt collection")
			return nil, err
		}
		return diff(watchlist, collection), nil

	}
	return watchlist, nil

}
