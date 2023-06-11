package cmd

import (
	"os"

	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/jokin1999/certark/tank"
	"github.com/spf13/cobra"
)

const (
	serviceConfigDir  = "/etc/certark"
	initLockFile      = ".lock"
	initLockFilePath  = serviceConfigDir + "/" + initLockFile
	serviceConfigFile = "config.yml"
	serviceConfigPath = serviceConfigDir + "/" + serviceConfigFile
	taskConfigDir     = serviceConfigDir + "/task"
	stateDir          = serviceConfigDir + "/state"
	acmeUserDir       = serviceConfigDir + "/user"
	certarkService    = "certark.service"
)

var rootCmd = &cobra.Command{}

func init() {
	rootCmd = &cobra.Command{
		Use:   "certark",
		Short: "CertArk is a certificate requestor based on lego.",
		Run: func(cmd *cobra.Command, args []string) {
			// show help information
			cmd.Help()
		},
	}
}

func Execute(version string) {
	// disable completion options
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// select running mode
	if version == "dev" {
		ark.Debug().Msg("Running in developing mode")
		tank.Save("MODE", certark.MODE_DEV)
	} else {
		ark.Debug().Msg("Running in product mode")
		tank.Save("MODE", certark.MODE_PROD)
	}

	if err := rootCmd.Execute(); err != nil {
		ark.Error().Err(err).Msg("Failed to run certark")
		os.Exit(1)
	}
}
