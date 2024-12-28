package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var defaultConfig = Config{
	API: APIConfig{
		TMDB: TMDBConfig{
			Language: "fr-FR",
		},
	},
	Scanner: ScannerConfig{
		MediaPath:       "./",
		Recursive:       true,
		IncludeNotFound: false,
	},
	Renamer: RenamerConfig{
		DryRun:     true,
		Type:       Movie,
		MaxResults: 5,
		QuickMode:  false,
		Patterns: PatternConfig{
			Movie:  "{name} - {year}{extension}",
			TVShow: "{name} - {season}x{episode}{extension}",
		},
	},
}

type MediaType string

const (
	Movie  MediaType = "movie"
	TvShow MediaType = "tvshow"
)

type Config struct {
	API     APIConfig     `yaml:"api"`
	Scanner ScannerConfig `yaml:"scanner"`
	Renamer RenamerConfig `yaml:"renamer"`
}

type APIConfig struct {
	TMDB TMDBConfig `yaml:"tmdb"`
}

type TMDBConfig struct {
	Key      string `yaml:"key"`
	Language string `yaml:"language"`
}

type ScannerConfig struct {
	MediaPath       string `yaml:"media_path"`
	Recursive       bool   `yaml:"recursive"`
	IncludeNotFound bool   `yaml:"include_not_found"`
}

type RenamerConfig struct {
	DryRun     bool          `yaml:"dry_run"`
	Type       MediaType     `yaml:"type"`
	Patterns   PatternConfig `yaml:"patterns"`
	MaxResults int           `yaml:"max_results"`
	QuickMode  bool          `yaml:"quick_mode"`
}

type PatternConfig struct {
	Movie  string `yaml:"movie"`
	TVShow string `yaml:"tvshow"`
}

// LoadConfig loads the configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config.yml"
	}

	f, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %w", err)
	}

	// Apply default values before validation
	cfg.applyDefaults()

	// Convert relative paths to absolute
	absPath, err := filepath.Abs(cfg.Scanner.MediaPath)
	if err != nil {
		return nil, fmt.Errorf("invalid media path: %w", err)
	}
	cfg.Scanner.MediaPath = absPath

	return &cfg, nil
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig(path string) error {
	cfg := defaultConfig

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	return nil
}
