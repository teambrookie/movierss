package worker

import (
	"log"
	"time"

	"github.com/teambrookie/movierss/dao"
	"github.com/teambrookie/movierss/torrent"
)

const apiRateLimit = 2 * time.Second

func Rarbg(in <-chan dao.Movie, out chan<- dao.Movie, config torrent.Config) {
	for movie := range in {
		time.Sleep(apiRateLimit)
		torrent, err := torrent.Search(movie.Ids.Imdb, config)
		if torrent.Download == "" {
			log.Printf("%s NOT FOUND", movie.Title)
			continue
		}
		if err != nil {
			log.Printf("Error processing %s : %s ...\n", movie.Title, err)
			continue
		}
		log.Printf(torrent.Title)
		movie.MagnetLink = torrent.Download
		movie.LastModified = time.Now()
		out <- movie
	}
}
