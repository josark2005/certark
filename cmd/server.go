package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var serverPort = ""

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start a CertArk server.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(serverPort)
		},
	}

	// Specify a server port
	// This flag has a higher priority than configuration and environmnet variables.
	serverCmd.Flags().StringVarP(&serverPort, "port", "p", "", "specify server listening tcp port")

	rootCmd.AddCommand(serverCmd)
}
