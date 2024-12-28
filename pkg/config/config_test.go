package config

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestApplyDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    Config
		expected Config
	}{
		{
			name:     "Empty config gets all defaults",
			input:    Config{},
			expected: defaultConfig,
		},
		{
			name: "Partial config keeps custom values",
			input: Config{
				API: APIConfig{
					TMDB: TMDBConfig{
						Language: "en-US",
					},
				},
			},
			expected: Config{
				API: APIConfig{
					TMDB: TMDBConfig{
						Language: "en-US",
					},
				},
				Scanner: defaultConfig.Scanner,
				Renamer: defaultConfig.Renamer,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.input
			cfg.applyDefaults()

			// Check specific fields that should have defaults
			if cfg.API.TMDB.Language != tt.expected.API.TMDB.Language {
				t.Errorf("Language = %v, want %v", cfg.API.TMDB.Language, tt.expected.API.TMDB.Language)
			}
			if cfg.Scanner.MediaPath != tt.expected.Scanner.MediaPath {
				t.Errorf("MediaPath = %v, want %v", cfg.Scanner.MediaPath, tt.expected.Scanner.MediaPath)
			}
			if cfg.Renamer.MaxResults != tt.expected.Renamer.MaxResults {
				t.Errorf("MaxResults = %v, want %v", cfg.Renamer.MaxResults, tt.expected.Renamer.MaxResults)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		shouldError bool
	}{
		{
			name: "Valid config",
			config: Config{
				API: APIConfig{
					TMDB: TMDBConfig{
						Key:      "valid-key",
						Language: "en-US",
					},
				},
				Scanner: defaultConfig.Scanner,
				Renamer: defaultConfig.Renamer,
			},
			shouldError: false,
		},
		{
			name: "Missing API key",
			config: Config{
				API: APIConfig{
					TMDB: TMDBConfig{
						Language: "en-US",
					},
				},
				Scanner: defaultConfig.Scanner,
				Renamer: defaultConfig.Renamer,
			},
			shouldError: true,
		},
		{
			name: "Invalid movie pattern",
			config: Config{
				API: APIConfig{
					TMDB: TMDBConfig{
						Key:      "valid-key",
						Language: "en-US",
					},
				},
				Scanner: defaultConfig.Scanner,
				Renamer: RenamerConfig{
					Patterns: PatternConfig{
						Movie: "invalid-pattern",
					},
				},
			},
			shouldError: true,
		},
		{
			name: "Invalid media type",
			config: Config{
				API: APIConfig{
					TMDB: TMDBConfig{
						Key: "valid-key",
					},
				},
				Renamer: RenamerConfig{
					Type: "invalid",
				},
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.shouldError {
				t.Errorf("validate() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	validConfig := []byte(`
api:
  tmdb:
    key: "test-key"
    language: "en-US"
scanner:
  media_path: "./test"
renamer:
  type: "movie"
`)

	if err := os.WriteFile(configPath, validConfig, 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	tests := []struct {
		name        string
		configPath  string
		shouldError bool
	}{
		{
			name:        "Valid config file",
			configPath:  configPath,
			shouldError: false,
		},
		{
			name:        "Non-existent file",
			configPath:  "non-existent.yml",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadConfig(tt.configPath)
			if (err != nil) != tt.shouldError {
				t.Errorf("LoadConfig() error = %v, shouldError %v", err, tt.shouldError)
			}
		})
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	if err := CreateDefaultConfig(configPath); err != nil {
		t.Fatalf("CreateDefaultConfig() error = %v", err)
	}

	// Read the created file content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read created config file: %v", err)
	}

	// Decode without validation
	var cfg Config
	if err := yaml.NewDecoder(bytes.NewReader(content)).Decode(&cfg); err != nil {
		t.Fatalf("Failed to decode created config: %v", err)
	}

	// Verify some key defaults
	if cfg.API.TMDB.Language != defaultConfig.API.TMDB.Language {
		t.Errorf("Default language = %v, want %v", cfg.API.TMDB.Language, defaultConfig.API.TMDB.Language)
	}
	if cfg.Renamer.MaxResults != defaultConfig.Renamer.MaxResults {
		t.Errorf("Default max results = %v, want %v", cfg.Renamer.MaxResults, defaultConfig.Renamer.MaxResults)
	}
	if cfg.Scanner.MediaPath != defaultConfig.Scanner.MediaPath {
		t.Errorf("Default media path = %v, want %v", cfg.Scanner.MediaPath, defaultConfig.Scanner.MediaPath)
	}
	if cfg.Renamer.Type != defaultConfig.Renamer.Type {
		t.Errorf("Default type = %v, want %v", cfg.Renamer.Type, defaultConfig.Renamer.Type)
	}
}
