package cmd

import (
	"fmt"
	"os"

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

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
