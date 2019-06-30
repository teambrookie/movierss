package torrent

import (
	torrentapi "github.com/qopher/go-torrentapi"
)

func excludeNoSeeder(torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	var movies torrentapi.TorrentResults
	for _, t := range torrents {
		if t.Seeders > 0 {
			movies = append(movies, t)
		}
	}
	return movies
}

func bestTorrent(torrents torrentapi.TorrentResults) torrentapi.TorrentResult {
	bt := torrentapi.TorrentResult{}
	for _, t := range torrents {
		if (bt == torrentapi.TorrentResult{}) {
			bt = t
			continue
		}
		if (t.Seeders / (1 + t.Leechers)) > (bt.Seeders / (1 + bt.Leechers)) {
			bt = t
		}
	}
	return bt
}

//Search is a function that search a movie on rarbg using an IMDB id
//by default it search the movie in category 44 also know as Serie/720p
func Search(movieIMBDID string, config Config) (torrentapi.TorrentResult, error) {
	api, err := torrentapi.New("Movierss")
	if err != nil {
		return torrentapi.TorrentResult{}, err
	}
	api.Format("json_extended")
	api.SearchIMDb(movieIMBDID)
	results, err := api.Search()
	if err != nil {
		return torrentapi.TorrentResult{}, err
	}

	if len(results) == 0 {
		return torrentapi.TorrentResult{}, nil
	}
	torrents := excludeNoSeeder(results)
	torrents = Filter(config.Categories, torrents)
	return bestTorrent(torrents), nil
}
