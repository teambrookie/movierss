package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/teambrookie/movierss/dao"
	"github.com/teambrookie/movierss/handlers"
	"github.com/teambrookie/movierss/torrent"
	"github.com/teambrookie/movierss/trakt"
)

const apiRateLimit = 2

func worker(jobs <-chan dao.Movie, store dao.MovieStore) {
	for movie := range jobs {
		time.Sleep(apiRateLimit * time.Second)
		log.Println("Processing : " + movie.Title)
		torrentLink, err := torrent.Search(movie.Ids.Imdb)
		log.Println("Result : " + torrentLink)
		if err != nil {
			log.Printf("Error processing %s : %s ...\n", movie.Title, err)
			continue
		}
		if torrentLink == "" {
			continue
		}
		movie.MagnetLink = torrentLink
		movie.LastModified = time.Now()
		err = store.UpdateMovie(movie)
		if err != nil {
			log.Printf("Error saving %s to DB ...\n", movie.Title)
		}
	}
}

func main() {
	var httpAddr = flag.String("http", "0.0.0.0:8000", "HTTP service address")
	var dbAddr = flag.String("db", "movierss.db", "DB address")
	flag.Parse()

	traktAPIKey := os.Getenv("TRAKT_KEY")
	if traktAPIKey == "" {
		log.Fatalln("TRAKT_KEY must be set in env")
	}

	movieProvider := trakt.Trakt{APIKey: traktAPIKey}

	log.Println("Starting server ...")
	log.Printf("HTTP service listening on %s", *httpAddr)
	log.Println("Connecting to db ...")

	//DB stuff
	store, err := dao.InitDB(*dbAddr)
	if err != nil {
		log.Fatalln("Error connecting to DB")
	}

	err = store.CreateBucket("movies")
	if err != nil {
		log.Fatalln("Error when creating bucket")
	}

	//Worker stuff
	log.Println("Starting worker ...")
	jobs := make(chan dao.Movie, 100)
	go worker(jobs, store)

	errChan := make(chan error, 10)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)
	mux.Handle("/movies", handlers.MovieHandler(store))
	mux.Handle("/refresh", handlers.RefreshHandler(store, movieProvider, jobs))
	mux.Handle("/rss", handlers.RSSHandler(store, movieProvider))

	httpServer := http.Server{}
	httpServer.Addr = *httpAddr
	httpServer.Handler = handlers.LoggingHandler(mux)

	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			httpServer.Shutdown(context.Background())
			os.Exit(0)
		}
	}
}
