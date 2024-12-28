package handlers

import (
	"context"

	"github.com/nouuu/gonamer/pkg/config"
)

// MediaHandler définit l'interface commune pour tous les handlers de médias
type MediaHandler interface {
	// Handle traite une suggestion de média
	Handle(ctx context.Context) error
}

// BaseHandler contient les éléments communs à tous les handlers
type BaseHandler struct {
	config     *config.Config
	DryRun     bool
	QuickMode  bool
	MaxResults int
}

// MediaSuggestion reste l'interface commune pour les suggestions
type MediaSuggestion interface {
	GetOriginalFilename() string
	GenerateFilename(pattern string) string
	HasSingleSuggestion() bool
	RenameFileManually(ctx context.Context) error
}

// NewBaseHandler crée un nouveau BaseHandler avec la configuration donnée
func NewBaseHandler(config *config.Config) BaseHandler {
	return BaseHandler{
		config:     config,
		DryRun:     config.Renamer.DryRun,
		QuickMode:  config.Renamer.QuickMode,
		MaxResults: config.Renamer.MaxResults,
	}
}
