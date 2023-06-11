package cmd

import (
	"fmt"
	"os"

	"github.com/jokin1999/certark/ark"
	"github.com/spf13/cobra"
)

func init() {
	// config main command
	var configCmd = cmdConfig()

	// config show
	configCmd.AddCommand(cmdConfigShow())

	rootCmd.AddCommand(configCmd)
}

// config command
func cmdConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "CertArk configurations",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
}

// config show
func cmdConfigShow() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show CertArk configuration",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			showConfig()
		},
	}
}

// show config
func showConfig() {
	configFile := serviceConfigPath

	// read file
	profileContent, err := os.ReadFile(configFile)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to read config file")
		return
	}

	fmt.Println(string(profileContent))
}

// config set
