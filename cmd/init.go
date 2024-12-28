package cmd

import (
	"github.com/nouuu/gonamer/cmd/cli/ui"
	"github.com/nouuu/gonamer/pkg/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize GoNamer configuration file",
	Long:  `Initialize GoNamer configuration file with default values.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if err := config.CreateDefaultConfig(cfgFile); err != nil {
			ui.ShowError(ctx, "Failed to create default configuration file %s: %v", cfgFile, err)
			return err
		}
		ui.ShowSuccess(ctx, "Default configuration file %s created successfully", cfgFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
