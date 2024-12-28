package ui

import (
	"context"

	"github.com/pterm/pterm"
)

func HandlePbStop(ctx context.Context, pb *pterm.ProgressbarPrinter) {
	if _, err := pb.Stop(); err != nil {
		ShowError(ctx, "Error stopping progress bar: %v", err)
	}
}
