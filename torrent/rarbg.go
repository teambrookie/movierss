package torrent

import (
	"strings"

	"log"

	torrentapi "github.com/qopher/go-torrentapi"
)

func filterMovies(torrents torrentapi.TorrentResults) string {
	var moviesextended torrentapi.TorrentResults
	// Search for extended version
	for _, t := range torrents {
		var filename = strings.ToLower(t.Filename)
		if strings.Contains(filename, "extended") {
			moviesextended = append(moviesextended, t)
		}
	}
	log.Println(torrents)
	var results torrentapi.TorrentResults
	results = filteraudioQuality("DTS-HD", moviesextended)
	//log.Printf("For quality %s the number of result if %d", "DTS-HD", len(results))
	if len(results) > 0 {
		return results[0].Download
	}
	results = filteraudioQuality("DTS-HD.MA.7.1", moviesextended)
	//log.Printf("For quality %s the number of result if %d", "DTS-HD.MA.7.1", len(results))
	if len(results) > 0 {
		return results[0].Download
	}
	results = filteraudioQuality("TrueHD.7.1Atmos", moviesextended)
	//log.Printf("For quality %s the number of result if %d", "TrueHD.7.1Atmos", len(results))
	if len(results) > 0 {
		return results[0].Download
	}
	results = filteraudioQuality("DTS", moviesextended)
	//log.Printf("For quality %s the number of result if %d", "DTS", len(results))
	if len(results) > 0 {
		return results[0].Download
	}

	results = filteraudioQuality("DTS-HD", torrents)
	//log.Printf("For quality %s the number of result if %d", "DTS-HD", len(results))
	if len(results) > 0 {
		return results[0].Download
	}
	results = filteraudioQuality("DTS-HD.MA.7.1", torrents)
	//log.Printf("For quality %s the number of result if %d", "DTS-HD.MA.7.1", len(results))
	if len(results) > 0 {
		return results[0].Download
	}
	results = filteraudioQuality("TrueHD.7.1Atmos", torrents)
	//log.Printf("For quality %s the number of result if %d", "TrueHD.7.1Atmos", len(results))
	if len(results) > 0 {
		return results[0].Download
	}
	results = filteraudioQuality("DTS", torrents)
	//log.Printf("For quality %s the number of result if %d", "DTS", len(results))
	if len(results) > 0 {
		return results[0].Download
	}

	return ""

}

func filteraudioQuality(quality string, torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	var movies torrentapi.TorrentResults
	for _, t := range torrents {
		var filename = strings.ToLower(t.Download)
		quality = strings.ToLower(quality)
		if strings.Contains(filename, quality) && t.Seeders > 0 {
			movies = append(movies, t)
		}
	}
	return movies
}

func Search(movieIMBDID, quality string) (string, error) {
	api, err := torrentapi.Init()
	if err != nil {
		return "", err
	}
	api.Format("json_extended")
	api.Category(44)
	api.SearchImDB(movieIMBDID)
	results, err := api.Search()
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", nil
	}
	return filterMovies(results), nil
}
