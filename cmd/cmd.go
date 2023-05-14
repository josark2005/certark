package cmd

import (
	"fmt"
	"os"

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
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
