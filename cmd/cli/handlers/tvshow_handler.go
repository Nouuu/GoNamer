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

type TvShowHandler struct {
	BaseHandler
	suggestions  mediarenamer.EpisodeSuggestions
	tvClient     mediadata.TvShowClient
	mediaRenamer *mediarenamer.MediaRenamer
	exitFunc     func() error
}

func NewTvShowHandler(
	base BaseHandler,
	suggestions mediarenamer.EpisodeSuggestions,
	tvClient mediadata.TvShowClient,
	mediaRenamer *mediarenamer.MediaRenamer,
	exitFunc func() error,
) *TvShowHandler {
	return &TvShowHandler{
		BaseHandler:  base,
		suggestions:  suggestions,
		tvClient:     tvClient,
		mediaRenamer: mediaRenamer,
		exitFunc:     exitFunc,
	}
}

func (h *TvShowHandler) Handle(ctx context.Context) error {
	if len(h.suggestions.SuggestedEpisodes) != 1 {
		return h.handleOptions(ctx)
	}

	if h.QuickMode {
		ui.ShowSuccess("Quick - Renaming episode %s", pterm.Yellow(h.suggestions.Episode.OriginalFilename))
		suggestion := h.suggestions.SuggestedEpisodes[0]
		return h.renameEpisode(ctx, h.suggestions, suggestion.TvShow, suggestion.Episode)
	}

	return h.handleOptions(ctx)
}

func (h *TvShowHandler) handleOptions(ctx context.Context) error {
	menuBuilder := ui.NewMenuBuilder()

	for _, episode := range h.suggestions.SuggestedEpisodes {
		episode := episode
		label := fmt.Sprintf("%s - %dx%02d - %s",
			episode.TvShow.Title,
			episode.Episode.SeasonNumber,
			episode.Episode.EpisodeNumber,
			episode.Episode.Name,
		)
		menuBuilder.AddOption(label, func() error {
			return h.renameEpisode(ctx, h.suggestions, episode.TvShow, episode.Episode)
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
			ui.ShowInfo("Skipping renaming of %s", pterm.Yellow(h.suggestions.Episode.OriginalFilename))
			return nil
		},
		func() error {
			return h.exitFunc()
		},
	)

	return menuBuilder.Build()
}

func (h *TvShowHandler) handleManualSearch(ctx context.Context) error {
	query, err := ui.PromptText(
		fmt.Sprintf("Search for '%s'", h.suggestions.Episode.OriginalFilename),
		h.suggestions.Episode.Name,
	)
	if err != nil {
		return err
	}

	tvShows, err := h.tvClient.SearchTvShow(ctx, query, 0, 1)
	if err != nil {
		return fmt.Errorf("error searching for tv show: %w", err)
	}

	h.suggestions.SuggestedEpisodes = make([]mediarenamer.SuggestedEpisode, 0, len(tvShows.TvShows))
	for _, tvShow := range tvShows.TvShows {
		episode, err := h.tvClient.GetEpisode(ctx, tvShow.ID, h.suggestions.Episode.Season, h.suggestions.Episode.Episode)
		if err != nil {
			ui.ShowError("Error getting episode: %v", err)
			continue
		}
		h.suggestions.SuggestedEpisodes = append(h.suggestions.SuggestedEpisodes, mediarenamer.SuggestedEpisode{
			TvShow:  tvShow,
			Episode: episode,
		})
	}

	if len(h.suggestions.SuggestedEpisodes) > h.config.Renamer.MaxResults {
		h.suggestions.SuggestedEpisodes = h.suggestions.SuggestedEpisodes[:h.config.Renamer.MaxResults]
	}

	return h.handleOptions(ctx)
}

func (h *TvShowHandler) handleManualRename(ctx context.Context) error {
	ui.ShowInfo("Renaming manually for %s", pterm.Yellow(h.suggestions.Episode.OriginalFilename))

	defaultValue := fmt.Sprintf("%s - %dx%02d",
		h.suggestions.Episode.Name,
		h.suggestions.Episode.Season,
		h.suggestions.Episode.Episode,
	)
	result, err := ui.PromptText("Enter new filename (without extension)", defaultValue)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.%s", result, h.suggestions.Episode.Extension)
	ui.ShowInfo("Renaming episode %s to %s",
		pterm.Yellow(h.suggestions.Episode.OriginalFilename),
		pterm.Yellow(filename),
	)

	return h.mediaRenamer.RenameFile(
		ctx,
		h.suggestions.Episode.FullPath,
		filepath.Join(filepath.Dir(h.suggestions.Episode.FullPath), filename),
		h.config.Renamer.DryRun,
	)
}

func (h *TvShowHandler) renameEpisode(
	ctx context.Context,
	suggestion mediarenamer.EpisodeSuggestions,
	tvShow mediadata.TvShow,
	episode mediadata.Episode,
) error {
	newFilename := mediarenamer.GenerateEpisodeFilename(h.config.Renamer.Patterns.TVShow, tvShow, episode, suggestion.Episode)
	if newFilename == suggestion.Episode.OriginalFilename {
		ui.ShowSuccess("Original filename is already correct for %s", pterm.Yellow(suggestion.Episode.OriginalFilename))
		return nil
	}

	ui.ShowInfo("Renaming episode %s to %s",
		pterm.Yellow(suggestion.Episode.OriginalFilename),
		pterm.Yellow(newFilename),
	)

	err := h.mediaRenamer.RenameEpisode(ctx, suggestion.Episode, tvShow, episode, h.config.Renamer.Patterns.TVShow, h.config.Renamer.DryRun)
	if err != nil {
		ui.ShowError("Error renaming episode: %v", err)
		return err
	}
	return nil
}
