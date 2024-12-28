package cli

import (
	"context"
	"errors"

	"github.com/nouuu/gonamer/cmd/cli/handlers"
	"github.com/nouuu/gonamer/cmd/cli/ui"
	"github.com/nouuu/gonamer/internal/mediadata"
	"github.com/nouuu/gonamer/internal/mediarenamer"
	"github.com/nouuu/gonamer/internal/mediascanner"
	"github.com/nouuu/gonamer/pkg/config"
	"github.com/pterm/pterm"
)

type Cli struct {
	config       *config.Config
	scanner      mediascanner.MediaScanner
	mediaRenamer *mediarenamer.MediaRenamer
	tvClient     mediadata.TvShowClient
	movieClient  mediadata.MovieClient
}

var ErrExit = errors.New("exit requested")

func NewCli(scanner mediascanner.MediaScanner, mediaRenamer *mediarenamer.MediaRenamer, movieClient mediadata.MovieClient, tvClient mediadata.TvShowClient, config *config.Config) *Cli {
	return &Cli{
		config:       config,
		scanner:      scanner,
		mediaRenamer: mediaRenamer,
		movieClient:  movieClient,
		tvClient:     tvClient,
	}
}

func (c *Cli) Run(ctx context.Context) error {
	switch c.config.Renamer.Type {
	case config.Movie:
		return c.processMovie(ctx)
	case config.TvShow:
		return c.processTvShow(ctx)
	}
	return nil
}

func (c *Cli) processMovie(ctx context.Context) error {
	movies, err := c.ScanMovies(ctx)
	if err != nil {
		ui.ShowError(ctx, "Error scanning movies: %v", err)
		return err
	}

	return c.processMoviesList(ctx, movies)
}

func (c *Cli) processMoviesList(ctx context.Context, movies []mediascanner.Movie) error {
	pb, _ := pterm.DefaultProgressbar.
		WithTotal(len(movies)).
		WithCurrent(1).
		WithTitle("Processing movies...").Start()
	defer ui.HandlePbStop(ctx, pb)

	for i, movie := range movies {
		ui.HandlePbStop(ctx, pb)

		// Affiche le titre du film en cours
		ui.ShowInfo(ctx, "Processing movie %d/%d: %s", i+1, len(movies), pterm.Yellow(movie.OriginalFilename))

		// Recherche des suggestions
		suggestions, err := c.mediaRenamer.SuggestMovies(ctx, movie, c.config.Renamer.MaxResults)
		if err != nil {
			suggestions = mediarenamer.MovieSuggestions{Movie: movie}
			ui.ShowError(ctx, "Error finding suggestions for %s: %v", movie.OriginalFilename, err)
		}

		// Création et exécution du handler
		handler := handlers.NewMovieHandler(
			handlers.NewBaseHandler(c.config),
			suggestions,
			c.movieClient,
			c.mediaRenamer,
			func() error { return ErrExit },
		)

		if err := handler.Handle(ctx); err != nil {
			if errors.Is(err, ErrExit) {
				return c.Exit()
			}
			ui.ShowError(ctx, "Error handling movie: %v", err)
		}
		pb, _ = pterm.DefaultProgressbar.
			WithTotal(len(movies)).
			WithCurrent(i + 2).
			WithTitle("Processing movies...").
			Start()
	}

	ui.ShowSuccess(ctx, "Finished processing movies")
	return nil
}

func (c *Cli) processTvShow(ctx context.Context) error {
	episodes, err := c.ScanTvEpisodes(ctx)
	if err != nil {
		ui.ShowError(ctx, "Error scanning tv shows: %v", err)
		return err
	}

	return c.processEpisodesList(ctx, episodes)
}

func (c *Cli) processEpisodesList(ctx context.Context, episodes []mediascanner.Episode) error {
	pb, _ := pterm.DefaultProgressbar.WithTotal(len(episodes)).WithCurrent(1).WithTitle("Processing episodes...").Start()
	defer ui.HandlePbStop(ctx, pb)

	for i, episode := range episodes {
		// Stop la barre avant d'afficher le menu
		ui.HandlePbStop(ctx, pb)

		// Affiche le titre de l'épisode en cours
		ui.ShowInfo(ctx, "Processing episode %d/%d: %s", i+1, len(episodes), pterm.Yellow(episode.OriginalFilename))

		// Recherche des suggestions
		suggestions, err := c.mediaRenamer.SuggestEpisodes(ctx, episode, c.config.Renamer.MaxResults)
		if err != nil {
			suggestions = mediarenamer.EpisodeSuggestions{Episode: episode}
			ui.ShowError(ctx, "Error finding suggestions for %s: %v", episode.OriginalFilename, err)
		}

		// Création et exécution du handler
		handler := handlers.NewTvShowHandler(
			handlers.NewBaseHandler(c.config),
			suggestions,
			c.tvClient,
			c.mediaRenamer,
			func() error { return ErrExit },
		)

		if err := handler.Handle(ctx); err != nil {
			if errors.Is(err, ErrExit) {
				return c.Exit()
			}
			ui.ShowError(ctx, "Error handling episode: %v", err)
		}

		// Redémarre la barre après le menu
		pb, _ = pterm.DefaultProgressbar.
			WithTotal(len(episodes)).
			WithTitle("Processing episodes...").
			WithCurrent(i + 2).
			Start()
	}

	ui.ShowSuccess(ctx, "Finished processing episodes")
	return nil
}

func (c *Cli) Exit() error {
	pterm.Info.Println("Exiting...")
	return nil
}
