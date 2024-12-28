package ui

import (
	"context"
	"fmt"

	"github.com/nouuu/gonamer/pkg/logger"
	"github.com/pterm/pterm"
)

// PromptText demande à l'utilisateur de saisir du texte
func PromptText(message, defaultValue string) (string, error) {
	result, err := pterm.DefaultInteractiveTextInput.
		WithDefaultValue(defaultValue).
		Show(message)

	if err != nil {
		return "", fmt.Errorf("error getting user input: %w", err)
	}

	return result, nil
}

// ShowInfo affiche un message d'information
func ShowInfo(ctx context.Context, message string, args ...interface{}) {
	log := logger.FromContext(ctx)
	pterm.Info.Printfln(message, args...)
	log.Infof(message, args...)
}

// ShowSuccess affiche un message de succès
func ShowSuccess(ctx context.Context, message string, args ...interface{}) {
	log := logger.FromContext(ctx)
	pterm.Success.Printfln(message, args...)
	log.Infof(message, args...)
}

// ShowWarning affiche un message d'avertissement
func ShowWarning(ctx context.Context, message string, args ...interface{}) {
	log := logger.FromContext(ctx)
	pterm.Warning.Printfln(message, args...)
	log.Warnf(message, args...)
}

// ShowError affiche un message d'erreur
func ShowError(ctx context.Context, message string, args ...interface{}) {
	log := logger.FromContext(ctx)
	pterm.Error.Printfln(message, args...)
	log.Errorf(message, args...)
}

func ShowWelcomeHeader() {
	pterm.DefaultHeader.Println("Media Renamer")
	pterm.Print("\n\n")
}
