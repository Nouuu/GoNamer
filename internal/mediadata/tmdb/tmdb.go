package tmdb

import (
	"github.com/cyruzin/golang-tmdb"
	"github.com/nouuu/mediatracker/internal/mediadata"
	"log/slog"
	"path"
	"strconv"
)

const tmdbImageBaseUrl = "https://image.tmdb.org/t/p/original"

type OptFunc func(opts *Opts)

type AllOpts struct {
	APIKey string
	Opts
}

type Opts struct {
	Lang  string
	Adult bool
}

func WithLang(lang string) OptFunc {
	return func(opts *Opts) {
		opts.Lang = lang
	}
}

func WithAdult(adult bool) OptFunc {
	return func(opts *Opts) {
		opts.Adult = adult
	}
}

func defaultOpts(apiKey string) AllOpts {
	return AllOpts{
		APIKey: apiKey,
		Opts: Opts{
			Lang:  "en-US",
			Adult: false,
		},
	}
}

type tmdbClient struct {
	client *tmdb.Client
	opts   AllOpts
}

func cfgMap(cfg AllOpts) map[string]string {
	return map[string]string{
		"language":      cfg.Lang,
		"include_adult": strconv.FormatBool(cfg.Adult),
	}
}

func NewMovieClient(APIKey string, opts ...OptFunc) mediadata.MovieClient {
	o := defaultOpts(APIKey)
	for _, optF := range opts {
		optF(&o.Opts)
	}

	client, err := tmdb.Init(o.APIKey)
	if err != nil {
		slog.Error("Failed to initialize TMDB client", slog.Any("error", err))
	}
	return &tmdbClient{client: client, opts: o}
}

func (t *tmdbClient) SearchMovie(query string, page int) (mediadata.MovieResults, error) {
	searchMovies, err := t.client.GetSearchMovies(query, cfgMap(t.opts))
	if err != nil {
		return mediadata.MovieResults{}, err
	}
	var movies []mediadata.Movie = buildMovieFromResult(searchMovies.SearchMoviesResults)
	return mediadata.MovieResults{
		Movies:         movies,
		Totals:         searchMovies.TotalResults,
		ResultsPerPage: 20,
	}, nil
}

func (t *tmdbClient) GetMovie(id int) (mediadata.Movie, error) {
	movieDetails, err := t.client.GetMovieDetails(id, cfgMap(t.opts))
	if err != nil {
		return mediadata.Movie{}, err
	}
	return buildMovie(movieDetails), nil
}

func buildMovieFromResult(result *tmdb.SearchMoviesResults) []mediadata.Movie {
	var movies []mediadata.Movie = make([]mediadata.Movie, len(result.Results))
	for i, movie := range result.Results {
		movies[i] = buildMovie(&tmdb.MovieDetails{
			ID:          movie.ID,
			Title:       movie.Title,
			Overview:    movie.Overview,
			ReleaseDate: movie.ReleaseDate,
		})
	}
	return movies
}

func buildMovie(movie *tmdb.MovieDetails) mediadata.Movie {
	return mediadata.Movie{
		ID:          strconv.FormatInt(movie.ID, 10),
		Title:       movie.Title,
		Overview:    movie.Overview,
		ReleaseDate: movie.ReleaseDate,
	}
}
