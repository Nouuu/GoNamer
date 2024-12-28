package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/nouuu/gonamer/cmd/cli"
	"github.com/nouuu/gonamer/cmd/cli/ui"
	"github.com/nouuu/gonamer/internal/cache"
	"github.com/nouuu/gonamer/internal/mediadata/tmdb"
	"github.com/nouuu/gonamer/internal/mediarenamer"
	"github.com/nouuu/gonamer/internal/mediascanner/filescanner"
	"github.com/nouuu/gonamer/pkg/config"
	"github.com/nouuu/gonamer/pkg/logger"
	"github.com/pterm/pterm"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()

	// Parse command line flags
	configPath := flag.String("config", "config.yml", "path to configuration file")
	createConfig := flag.Bool("init", false, "create default configuration file")
	flag.Parse()

	// Create default config if requested
	if *createConfig {
		if err := config.CreateDefaultConfig(*configPath); err != nil {
			ui.ShowError("Failed to create default configuration %v", err)
			os.Exit(1)
		}
		ui.ShowSuccess("Default configuration file created at %s", *configPath)
		os.Exit(0)
	}

	// Initialize logger
	initLogger()

	// Start CLI with new configuration
	if err := startCli(ctx, *configPath); err != nil {
		fmt.Printf("Error starting application: %v\n", err)
		os.Exit(1)
	}
}

func initLogger() {
	logger.SetLoggerLevel(zapcore.InfoLevel)
	logfile, err := os.OpenFile("mediatracker.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		ui.ShowError("Error opening log file: %v", err)
		os.Exit(1)
	}

	logger.SetLoggerOutput(zapcore.WriteSyncer(logfile))
}

func startCli(ctx context.Context, configPath string) error {
	log := logger.FromContext(ctx)

	pterm.DefaultHeader.Println("Media Renamer")
	pterm.Print("\n\n")

	pterm.Info.Printfln("Loading configuration from %s...\n", configPath)

	// Load configuration from file
	conf, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if conf.Renamer.DryRun {
		pterm.Info.Println("Dry run enabled")
	} else {
		pterm.Warning.Println("Dry run disabled")
	}

	cacheClient, err := cache.NewGoCache(ctx)
	if err != nil {
		pterm.Error.Println(pterm.Error.Sprintf("Error creating cache client: %v", err))
		log.Fatalf("Error creating cache client: %v", err)
	}

	scanner := filescanner.New()
	movieClient, err := tmdb.NewMovieClient(conf.API.TMDB.Key, cacheClient, tmdb.WithLang(conf.API.TMDB.Language))
	if err != nil {
		pterm.Error.Printfln("Error creating movie client: %v", err)
		log.Fatalf("Error creating movie client: %v", err)
	}

	tvShowClient, err := tmdb.NewTvShowClient(conf.API.TMDB.Key, cacheClient, tmdb.WithLang(conf.API.TMDB.Language))
	if err != nil {
		pterm.Error.Println(pterm.Error.Sprintf("Error creating tv show client: %v", err))
		log.Fatalf("Error creating tv show client: %v", err)
	}

	mediaRenamer := mediarenamer.NewMediaRenamer(movieClient, tvShowClient)

	newCli := cli.NewCli(scanner, mediaRenamer, movieClient, tvShowClient, conf)

	newCli.Run(ctx)
	return nil
}
