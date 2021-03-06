package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"net/http"

	"github.com/teambrookie/movierss/dao"
	"github.com/teambrookie/movierss/handlers"
	"github.com/teambrookie/movierss/torrent"
	"github.com/teambrookie/movierss/trakt"
	"github.com/teambrookie/movierss/worker"
)

func loadConfig() (torrent.Config, error) {
	configURL := os.Getenv("SHOWRSS_CONFIGURL")
	config := torrent.Config{}
	if configURL != "" {
		resp, err := http.Get(configURL)
		if err != nil {
			return torrent.Config{}, err
		}
		err = json.NewDecoder(resp.Body).Decode(&config)
		if err != nil {
			return torrent.Config{}, err
		}
		return config, nil
	}
	configPath := os.Getenv("SHOWRSS_CONFIGPATH")
	if configPath == "" {
		configPath = "./config.json"
	}
	data, _ := os.Open(configPath)
	err := json.NewDecoder(bufio.NewReader(data)).Decode(&config)
	if err != nil {
		return torrent.Config{}, err
	}
	return config, nil

}

func main() {
	var httpAddr = flag.String("http", "0.0.0.0:8000", "HTTP service address")
	var dbAddr = flag.String("db", "movierss.db", "DB address")
	flag.Parse()

	traktAPIKey := os.Getenv("TRAKT_KEY")
	if traktAPIKey == "" {
		log.Fatalln("TRAKT_KEY must be set in env")
	}

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading the config: %s", err)
	}

	movieProvider := trakt.Trakt{APIKey: traktAPIKey}

	fmt.Println("Starting server ...")
	fmt.Printf("HTTP service listening on %s\n", *httpAddr)
	fmt.Println("Connecting to db ...")

	//DB stuff
	store, err := dao.InitDB(*dbAddr)
	if err != nil {
		log.Fatalln("Error connecting to DB")
	}

	err = store.CreateBucket("movies")
	if err != nil {
		log.Fatalln("Error when creating bucket")
	}

	//searchTorrentWorker stuff

	in := make(chan dao.Movie, 100)
	out := make(chan dao.Movie, 100)

	fmt.Println("Starting workers ...")
	go worker.DB(out, store)
	go worker.Rarbg(in, out, config)
	go worker.Cleanup(in, store)

	errChan := make(chan error, 10)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)
	mux.Handle("/movies", handlers.MovieHandler(store))
	mux.Handle("/refresh", handlers.RefreshHandler(store, movieProvider, in))
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
