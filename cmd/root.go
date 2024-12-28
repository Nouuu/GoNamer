package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	// Global flags
	cfgFile             string
	cfgFileFlag         = "config"
	cfgFileShort        = "c"
	dryRun              bool
	dryRunFlag          = "dry-run"
	recursive           bool
	recursiveFlag       = "recursive"
	recursiveShort      = "r"
	maxResults          int
	maxResultsFlag      = "max-results"
	maxResultsShort     = "m"
	quickMode           bool
	quickModeFlag       = "quick"
	quickModeShort      = "q"
	mediaType           string
	mediaTypeFlag       = "type"
	mediaTypeShort      = "t"
	includeNotFound     bool
	includeNotFoundFlag = "include-not-found"
	moviePattern        string
	moviePatternFlag    = "movie-pattern"
	tvshowPattern       string
	tvshowPatternFlag   = "tvshow-pattern"
	language            string
	languageFlag        = "language"
	languageShort       = "l"
)

var rootCmd = &cobra.Command{
	Use:     "gonamer",
	Short:   "GoNamer - Rename your media files using TMDB metadata",
	Long:    `GoNamer is a powerful media file renaming tool that uses the TMDB API to automatically organize and rename your movie and TV show files.`,
	Version: fmt.Sprintf("%s (commit: %s) built at %s", version, commit, date),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Initialize global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, cfgFileFlag, cfgFileShort, "config.yml", "config file (default is ./config.yml)")

	// Scanner flags
	rootCmd.PersistentFlags().BoolVarP(&recursive, recursiveFlag, recursiveShort, true, "scan directories recursively")
	rootCmd.PersistentFlags().BoolVar(&includeNotFound, includeNotFoundFlag, false, "include files without matches")

	// Renamer flags
	rootCmd.PersistentFlags().BoolVar(&dryRun, dryRunFlag, true, "simulate renaming without actual changes")
	rootCmd.PersistentFlags().StringVarP(&mediaType, mediaTypeFlag, mediaTypeShort, "movie", "media type (movie or tvshow)")
	rootCmd.PersistentFlags().IntVarP(&maxResults, maxResultsFlag, maxResultsShort, 5, "maximum number of suggestions")
	rootCmd.PersistentFlags().BoolVarP(&quickMode, quickModeFlag, quickModeShort, false, "quick mode without confirmation")

	// Pattern flags
	rootCmd.PersistentFlags().StringVar(&moviePattern, moviePatternFlag, "{name} - {year}{extension}", "movie renaming pattern")
	rootCmd.PersistentFlags().StringVar(&tvshowPattern, tvshowPatternFlag, "{name} - {season}x{episode}{extension}", "TV show renaming pattern")

	// API flags
	rootCmd.PersistentFlags().StringVarP(&language, languageFlag, languageShort, "en-US", "preferred language for TMDB API")

	rootCmd.SetVersionTemplate("GoNamer {{.Version}}\n")
}
