package cmd

import (
	"os"

	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/tank"
	"github.com/spf13/cobra"
)

const (
	serviceConfigDir  = "/etc/certark"
	initLockFile      = ".lock"
	initLockFilePath  = serviceConfigDir + "/" + initLockFile
	serviceConfigFile = "config.yml"
	serviceConfigPath = serviceConfigDir + "/" + serviceConfigFile
	domainConfigDir   = serviceConfigDir + "/domain"
	acmeUserDir       = serviceConfigDir + "/user"
	certarkService    = "certark.service"
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

	// select running mode
	if version == "dev" {
		tank.Save("MODE", tank.MODE_DEV)
	} else {
		tank.Save("MODE", tank.MODE_PROD)
	}

	if err := rootCmd.Execute(); err != nil {
		ark.Error().Err(err).Msg("Failed to run certark")
		os.Exit(1)
	}
}
