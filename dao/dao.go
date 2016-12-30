package dao

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

type Movie struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	Ids   struct {
		Trakt int    `json:"trakt"`
		Slug  string `json:"slug"`
		Imdb  string `json:"imdb"`
		Tmdb  int    `json:"tmdb"`
	} `json:"ids"`
	MagnetLink   string    `json:"magnet_link"`
	LastModified time.Time `json:"last_modified"`
}

type MovieStore interface {
	GetMovie(string) (Movie, error)
	AddMovie(Movie) error
	UpdateMovie(Movie) error
	DeleteMovie(string) error
	GetAllMovies() ([]Movie, error)
	GetAllNotFoundMovies() ([]Movie, error)
}

type BoltMovieStore struct {
	db *bolt.DB
}

func InitDB(dbName string) (*BoltMovieStore, error) {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &BoltMovieStore{db}, nil
}

func (store *BoltMovieStore) CreateBucket(bucketName string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	return err
}

func (store *BoltMovieStore) GetMovie(name string) (Movie, error) {
	var movie Movie
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("movies"))
		v := b.Get([]byte(name))
		json.Unmarshal(v, &movie)
		return nil
	})
	return movie, err

}

func (store *BoltMovieStore) AddMovie(mov Movie) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("movies"))
		if v := b.Get([]byte(mov.Title)); v != nil {
			return nil
		}
		encoded, err := json.Marshal(mov)
		if err != nil {
			return err
		}

		return b.Put([]byte(mov.Title), encoded)
	})
	return err
}

func (store *BoltMovieStore) UpdateMovie(mov Movie) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		encoded, err := json.Marshal(mov)
		if err != nil {
			return err
		}
		b := tx.Bucket([]byte("movies"))
		return b.Put([]byte(mov.Title), encoded)
	})
	return err
}

func (store *BoltMovieStore) DeleteMovie(name string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("movies"))
		return b.Delete([]byte(name))
	})
	return err
}

func (store *BoltMovieStore) GetAllMovie() ([]Movie, error) {
	var movies []Movie
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("movies"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var movie Movie
			json.Unmarshal(v, &movie)
			movies = append(movies, movie)
		}
		return nil
	})
	return movies, err
}

func (store *BoltMovieStore) GetAllNotFoundMovie() ([]Movie, error) {
	var movies []Movie
	err := store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("movies"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var movie Movie
			json.Unmarshal(v, &movie)
			if movie.MagnetLink == "" {
				movies = append(movies, movie)
			}

		}
		return nil
	})
	return movies, err
}
