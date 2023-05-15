package cmd

import (
	"fmt"

	"github.com/jokin1999/certark/ark"
	"github.com/spf13/cobra"
)

func init() {
	var serverPort = ""

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start a CertArk server",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			fmt.Println(serverPort)
		},
	}

	// Specify a server port
	// This flag has a higher priority than configuration and environmnet variables.
	serverCmd.Flags().StringVarP(&serverPort, "port", "p", "", "specify server listening tcp port")

	rootCmd.AddCommand(serverCmd)
}
