package ui

import (
	"context"

	"github.com/pterm/pterm"
)

func HandleSpinnerStop(ctx context.Context, spinner *pterm.SpinnerPrinter) {
	if err := spinner.Stop(); err != nil {
		ShowError(ctx, "Error stopping spinner: %v", err)
	}
}
