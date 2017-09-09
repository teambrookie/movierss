package worker

import (
	"log"
	"time"

	"github.com/teambrookie/movierss/dao"
	"github.com/teambrookie/movierss/torrent"
)

const apiRateLimit = 2 * time.Second

func Rarbg(in <-chan dao.Movie, out chan<- dao.Movie) {
	for movie := range in {
		time.Sleep(apiRateLimit)
		log.Println("Processing : " + movie.Title)
		torrentLink, err := torrent.Search(movie.Ids.Imdb)
		log.Println("Result : " + torrentLink)
		if err != nil {
			log.Printf("Error processing %s : %s ...\n", movie.Title, err)
			continue
		}
		movie.MagnetLink = torrentLink
		movie.LastModified = time.Now()
		out <- movie
	}
}
