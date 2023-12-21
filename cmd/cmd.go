package cmd

import (
	"os"

	"github.com/josark2005/certark/ark"
	"github.com/josark2005/certark/certark"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "certark",
	Short: "CertArk is a certificate requestor based on lego.",
	Run: func(cmd *cobra.Command, args []string) {
		// show help information
		cmd.Help()
	},
}

func Execute(version string) {
	// disable completion options
	// rootCmd.CompletionOptions.DisableDefaultCmd = true

	// select running mode
	if certark.CurrentConfig.Mode == "prod" {
		certark.CurrentConfig.Mode = certark.MODE_PROD
	} else {
		ark.Debug().Msg("Running in developing mode")
		certark.CurrentConfig.Mode = certark.MODE_DEV
	}

	if err := rootCmd.Execute(); err != nil {
		ark.Error().Err(err).Msg("Failed to run certark")
		os.Exit(1)
	}
}
