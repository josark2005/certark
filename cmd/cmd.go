package cmd

import (
	"os"

	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func init() {
	rootCmd = &cobra.Command{
		Use:   "certark",
		Short: "CertArk is a certificate requestor based on lego.",
		Run: func(cmd *cobra.Command, args []string) {
			// show help information
			cmd.Help()
			println(args)
		},
	}
}

func Execute(version string) {
	// disable completion options
	// rootCmd.CompletionOptions.DisableDefaultCmd = true

	// select running mode
	if certark.CurrentConfig.Mode == "prod" {
		ark.Debug().Msg("Running in product mode")
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
