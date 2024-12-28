package config

import (
	"fmt"
	"strings"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors holds multiple validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var errStrings []string
	for _, err := range e {
		errStrings = append(errStrings, err.Error())
	}
	return "validation errors:\n- " + strings.Join(errStrings, "\n- ")
}

// applyDefaults applies default values to empty fields in the configuration
func (c *Config) applyDefaults() {
	if c.API.TMDB.Language == "" {
		c.API.TMDB.Language = defaultConfig.API.TMDB.Language
	}

	if c.Scanner.MediaPath == "" {
		c.Scanner.MediaPath = defaultConfig.Scanner.MediaPath
	}

	if c.Renamer.MaxResults <= 0 {
		c.Renamer.MaxResults = defaultConfig.Renamer.MaxResults
	}

	if c.Renamer.Patterns.Movie == "" {
		c.Renamer.Patterns.Movie = defaultConfig.Renamer.Patterns.Movie
	}

	if c.Renamer.Patterns.TVShow == "" {
		c.Renamer.Patterns.TVShow = defaultConfig.Renamer.Patterns.TVShow
	}
}

// validate performs comprehensive validation of the configuration
func (c *Config) validate() error {
	var errs ValidationErrors

	// Validate required fields
	if c.API.TMDB.Key == "" {
		errs = append(errs, ValidationError{
			Field:   "api.tmdb.key",
			Message: "TMDB API key is required",
		})
	}

	// Validate patterns
	if !hasValidPatternVariables(c.Renamer.Patterns.Movie, []string{"{name}", "{year}", "{extension}"}) {
		errs = append(errs, ValidationError{
			Field:   "renamer.patterns.movie",
			Message: "movie pattern must contain {name}, {year}, and {extension}",
		})
	}

	if !hasValidPatternVariables(c.Renamer.Patterns.TVShow, []string{"{name}", "{season}", "{episode}", "{extension}"}) {
		errs = append(errs, ValidationError{
			Field:   "renamer.patterns.tvshow",
			Message: "tv show pattern must contain {name}, {season}, {episode}, and {extension}",
		})
	}

	// Validate media type
	if !isValidMediaType(c.Renamer.Type) {
		errs = append(errs, ValidationError{
			Field:   "renamer.type",
			Message: "invalid media type, must be 'movie' or 'tvshow'",
		})
	}

	// Validate numeric values
	if c.Renamer.MaxResults < 1 {
		errs = append(errs, ValidationError{
			Field:   "renamer.max_results",
			Message: "max_results must be greater than 0",
		})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func hasValidPatternVariables(pattern string, requiredVars []string) bool {
	for _, v := range requiredVars {
		if !strings.Contains(pattern, v) {
			return false
		}
	}
	return true
}

func isValidMediaType(t MediaType) bool {
	return t == Movie || t == TvShow
}
