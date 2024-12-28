package cmd

import (
	"context"
	"os"

	"github.com/nouuu/gonamer/cmd/cli"
	"github.com/nouuu/gonamer/cmd/cli/ui"
	"github.com/nouuu/gonamer/internal/cache"
	"github.com/nouuu/gonamer/internal/mediadata/tmdb"
	"github.com/nouuu/gonamer/internal/mediarenamer"
	"github.com/nouuu/gonamer/internal/mediascanner/filescanner"
	"github.com/nouuu/gonamer/pkg/config"
	"github.com/nouuu/gonamer/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
)

var renameCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "rename [path]",
	Short: "Rename media files using TMDB metadata.",
	Long: `Rename media files in the specified path using TMDB metadata.
If no path is specified, the current directory will be used.`,
	RunE: runRename,
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	if err := initLogger(ctx); err != nil {
		return err
	}

	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	ui.ShowInfo(ctx, "Loading configuration file from '%s'\n", cfgFile)

	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		ui.ShowError(ctx, "Failed to load configuration file '%s': %v", cfgFile, err)
		return err
	}

	// Override with flags
	if cmd.Flags().Changed(dryRunFlag) {
		cfg.Renamer.DryRun = dryRun
	}
	if cmd.Flags().Changed(recursiveFlag) {
		cfg.Scanner.Recursive = recursive
	}
	if cmd.Flags().Changed(maxResultsFlag) {
		cfg.Renamer.MaxResults = maxResults
	}
	if cmd.Flags().Changed(quickModeFlag) {
		cfg.Renamer.QuickMode = quickMode
	}
	if cmd.Flags().Changed(mediaTypeFlag) {
		cfg.Renamer.Type = config.MediaType(mediaType)
	}
	if cmd.Flags().Changed(moviePatternFlag) {
		cfg.Renamer.Patterns.Movie = moviePattern
	}
	if cmd.Flags().Changed(tvshowPatternFlag) {
		cfg.Renamer.Patterns.TVShow = tvshowPattern
	}
	if cmd.Flags().Changed(languageFlag) {
		cfg.API.TMDB.Language = language
	}
	if cmd.Flags().Changed(includeNotFoundFlag) {
		cfg.Scanner.IncludeNotFound = includeNotFound
	}

	if err := cfg.Validate(); err != nil {
		ui.ShowError(ctx, "Invalid configuration: %v", err)
		return err
	}

	return startCli(ctx, cfg, path)
}

func initLogger(ctx context.Context) error {
	logger.SetLoggerLevel(zapcore.InfoLevel)
	logfile, err := os.OpenFile("mediatracker.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		ui.ShowError(ctx, "Error opening log file: %v", err)
		return err
	}

	logger.SetLoggerOutput(zapcore.WriteSyncer(logfile))
	return nil
}

func startCli(ctx context.Context, conf *config.Config, mediaPath string) error {
	ui.ShowWelcomeHeader()

	if mediaPath != "." {
		ui.ShowInfo(ctx, "Using media path '%s' instead of the one in the configuration file", mediaPath)
		conf.Scanner.MediaPath = mediaPath
	}

	if conf.Renamer.DryRun {
		ui.ShowSuccess(ctx, "Dry run mode enabled, no files will be renamed")
	} else {
		ui.ShowWarning(ctx, "Dry run mode disabled, files will be renamed")
	}

	cacheClient, err := cache.NewGoCache(ctx)
	if err != nil {
		ui.ShowError(ctx, "Error creating cache client: %v", err)
		return err
	}

	scanner := filescanner.New()
	movieClient, err := tmdb.NewMovieClient(conf.API.TMDB.Key, cacheClient, tmdb.WithLang(conf.API.TMDB.Language))
	if err != nil {
		ui.ShowError(ctx, "Error creating movie client: %v", err)
		return err
	}

	tvShowClient, err := tmdb.NewTvShowClient(conf.API.TMDB.Key, cacheClient, tmdb.WithLang(conf.API.TMDB.Language))
	if err != nil {
		ui.ShowError(ctx, "Error creating tv show client: %v", err)
		return err
	}

	mediaRenamer := mediarenamer.NewMediaRenamer(movieClient, tvShowClient)

	newCli := cli.NewCli(scanner, mediaRenamer, movieClient, tvShowClient, conf)

	return newCli.Run(ctx)

}
