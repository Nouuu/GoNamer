package main

import (
	"context"

	"github.com/nouuu/mediatracker/conf"
	"github.com/nouuu/mediatracker/internal/mediadata/tmdb"
	"github.com/nouuu/mediatracker/internal/mediarenamer"
	"github.com/nouuu/mediatracker/internal/mediascanner"
	"github.com/nouuu/mediatracker/internal/mediascanner/filescanner"
	"github.com/nouuu/mediatracker/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger.SetLoggerLevel(zapcore.ErrorLevel)
	ctx := context.Background()
	log := logger.FromContext(ctx)

	config := conf.LoadConfig()

	scanner := filescanner.New()
	movieClient, err := tmdb.NewMovieClient(config.TMDBAPIKey, tmdb.WithLang("fr-FR"))
	if err != nil {
		log.Fatalf("Error creating movie client: %v", err)
	}
	mediaRenamer := mediarenamer.NewMediaRenamer(movieClient)

	movies, err := scanner.ScanMovies(ctx, "/mnt/nfs/Download/direct_download/film", mediascanner.ScanMoviesOptions{Recursively: true})
	if err != nil {
		log.Fatalf("Error scanning movies: %v", err)
	}
	mediaRenamer.FindSuggestions(ctx, movies)
	if err != nil {
		log.Fatalf("Error renaming movies: %v", err)
	}
}

func findSuggestionCallback(suggestion mediarenamer.MovieSuggestions, err error) {
	if err != nil {
		// Handle error
	}
	// Do something with suggestion
}
