package tmdb

import (
	"github.com/cyruzin/golang-tmdb"
	"github.com/nouuu/mediatracker/internal/mediadata"
	"log/slog"
	"strconv"
)

func NewTvShowClient(APIKey string, opts ...OptFunc) mediadata.TvShowClient {
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

func (t *tmdbClient) SearchTvShow(query string, page int) (mediadata.TvShowResults, error) {
	searchTvShows, err := t.client.GetSearchTVShow(query, cfgMap(t.opts, map[string]string{
		"page": strconv.Itoa(page),
	}))
	if err != nil {
		return mediadata.TvShowResults{}, err
	}
	var tvShows []mediadata.TvShow = buildTvShowFromResult(searchTvShows.SearchTVShowsResults)
	return mediadata.TvShowResults{
		TvShows:        tvShows,
		Totals:         searchTvShows.TotalResults,
		ResultsPerPage: 20,
	}, nil
}

func (t *tmdbClient) GetTvShow(id string) (mediadata.TvShow, error) {
	var idInt int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return mediadata.TvShow{}, err
	}
	tvShowDetails, err := t.client.GetTVDetails(
		idInt,
		cfgMap(t.opts, map[string]string{
			"append_to_response": "credits",
		}),
	)
	if err != nil {
		return mediadata.TvShow{}, err
	}
	return buildTvShow(tvShowDetails), nil
}

func buildTvShow(tvShow *tmdb.TVDetails) mediadata.TvShow {
	return mediadata.TvShow{
		ID:          strconv.FormatInt(tvShow.ID, 10),
		Title:       tvShow.Name,
		Overview:    tvShow.Overview,
		FistAirDate: tvShow.FirstAirDate,
		PosterURL:   tmdbImageBaseUrl + tvShow.PosterPath,
		Rating:      tvShow.VoteAverage,
		RatingCount: tvShow.VoteCount,
	}
}

func buildTvShowFromResult(result *tmdb.SearchTVShowsResults) []mediadata.TvShow {
	var tvShows = make([]mediadata.TvShow, len(result.Results))
	for i, tvShow := range result.Results {
		tvShows[i] = buildTvShow(&tmdb.TVDetails{
			ID:           tvShow.ID,
			Name:         tvShow.Name,
			Overview:     tvShow.Overview,
			FirstAirDate: tvShow.FirstAirDate,
			PosterPath:   tvShow.PosterPath,
			VoteAverage:  tvShow.VoteAverage,
			VoteCount:    tvShow.VoteCount,
		})
	}
	return tvShows
}