package cmd

import (
	"os"

	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
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

	// load config
	config, err := ReadConfig()
	if err != nil {
		ark.Warn().Err(err).Msg("Load CertArk config failed, may fallback to default")
	}
	certark.CurrentConfig = config
}

func Execute(version string) {
	// disable completion options
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// select running mode
	if certark.CurrentConfig.Mode == "dev" {
		ark.Debug().Msg("Running in developing mode")
		certark.CurrentConfig.Mode = certark.MODE_DEV
	} else {
		ark.Debug().Msg("Running in product mode")
		certark.CurrentConfig.Mode = certark.MODE_PROD
	}

	if err := rootCmd.Execute(); err != nil {
		ark.Error().Err(err).Msg("Failed to run certark")
		os.Exit(1)
	}
}
