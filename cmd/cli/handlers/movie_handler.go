package handlers

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/nouuu/gonamer/cmd/cli/ui"
	"github.com/nouuu/gonamer/internal/mediadata"
	"github.com/nouuu/gonamer/internal/mediarenamer"
	"github.com/pterm/pterm"
)

type MovieHandler struct {
	BaseHandler
	suggestion   mediarenamer.MovieSuggestions
	movieClient  mediadata.MovieClient
	mediaRenamer *mediarenamer.MediaRenamer
	exitFunc     func() error
}

func NewMovieHandler(
	base BaseHandler,
	suggestion mediarenamer.MovieSuggestions,
	movieClient mediadata.MovieClient,
	mediaRenamer *mediarenamer.MediaRenamer,
	exitFunc func() error,
) *MovieHandler {
	return &MovieHandler{
		BaseHandler:  base,
		suggestion:   suggestion,
		movieClient:  movieClient,
		mediaRenamer: mediaRenamer,
		exitFunc:     exitFunc,
	}
}
func (h *MovieHandler) Handle(ctx context.Context) error {
	if len(h.suggestion.SuggestedMovies) != 1 {
		return h.handleOptions(ctx)
	}

	if h.QuickMode {
		ui.ShowSuccess("Quick - Renaming movie %s", pterm.Yellow(h.suggestion.Movie.OriginalFilename))
		return h.renameMovie(ctx, h.suggestion, h.suggestion.SuggestedMovies[0])
	}

	return h.handleOptions(ctx)
}

func (h *MovieHandler) handleOptions(ctx context.Context) error {
	menuBuilder := ui.NewMenuBuilder()

	for _, movie := range h.suggestion.SuggestedMovies {
		movie := movie
		label := fmt.Sprintf("%s (%s)", movie.Title, movie.Year)
		menuBuilder.AddOption(label, func() error {
			return h.renameMovie(ctx, h.suggestion, movie)
		})
	}

	menuBuilder.AddOption("Search Manually", func() error {
		return h.handleManualSearch(ctx)
	})

	menuBuilder.AddOption("Rename Manually", func() error {
		return h.handleManualRename(ctx)
	})

	menuBuilder.AddStandardOptions(
		func() error {
			ui.ShowInfo("Skipping renaming of %s", pterm.Yellow(h.suggestion.Movie.OriginalFilename))
			return nil
		},
		func() error {
			return h.exitFunc()
		},
	)

	return menuBuilder.Build()
}

func (h *MovieHandler) handleManualSearch(ctx context.Context) error {
	query, err := ui.PromptText(
		fmt.Sprintf("Search for '%s'", h.suggestion.Movie.OriginalFilename),
		h.suggestion.Movie.Name,
	)
	if err != nil {
		return err
	}

	movies, err := h.movieClient.SearchMovie(ctx, query, 0, 1)
	if err != nil {
		return fmt.Errorf("error searching for movie: %w", err)
	}

	h.suggestion.SuggestedMovies = movies.Movies
	if len(h.suggestion.SuggestedMovies) > h.config.Renamer.MaxResults {
		h.suggestion.SuggestedMovies = h.suggestion.SuggestedMovies[:h.config.Renamer.MaxResults]
	}

	return h.handleOptions(ctx)
}

func (h *MovieHandler) handleManualRename(ctx context.Context) error {
	ui.ShowInfo("Renaming manually for %s", pterm.Yellow(h.suggestion.Movie.OriginalFilename))

	defaultValue := fmt.Sprintf("%s (%d)", h.suggestion.Movie.Name, h.suggestion.Movie.Year)
	result, err := ui.PromptText("Enter new filename (without extension)", defaultValue)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.%s", result, h.suggestion.Movie.Extension)
	ui.ShowInfo("Renaming movie %s to %s", pterm.Yellow(h.suggestion.Movie.OriginalFilename), pterm.Yellow(filename))

	return h.mediaRenamer.RenameFile(
		ctx,
		h.suggestion.Movie.FullPath,
		filepath.Join(filepath.Dir(h.suggestion.Movie.FullPath), filename),
		h.config.Renamer.DryRun,
	)
}

func (h *MovieHandler) renameMovie(ctx context.Context, suggestion mediarenamer.MovieSuggestions, movie mediadata.Movie) error {
	newFilename := mediarenamer.GenerateMovieFilename(h.config.Renamer.Patterns.Movie, movie, suggestion.Movie)
	if newFilename == suggestion.Movie.OriginalFilename {
		ui.ShowSuccess("Original filename is already correct for %s", pterm.Yellow(h.suggestion.Movie.OriginalFilename))
		return nil
	}

	pterm.Info.Println("Renaming movie ", pterm.Yellow(suggestion.Movie.OriginalFilename), "to", pterm.Yellow(newFilename))
	err := h.mediaRenamer.RenameMovie(ctx, suggestion.Movie, movie, h.config.Renamer.Patterns.Movie, h.config.Renamer.DryRun)
	if err != nil {
		pterm.Error.Println(pterm.Sprintf("Error renaming movie: %v", err))
		return err
	}
	return nil
}
