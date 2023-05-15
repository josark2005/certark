package cmd

import (
	"fmt"
	"os"

	"github.com/jokin1999/certark/ark"
	"github.com/spf13/cobra"
)

func init() {
	// CertArk server link
	var serverLink = ""

	var clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Start a CertArk client",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			fmt.Println("Client is developing")
			os.Exit(0)
		},
	}

	// Specify a server port
	// This flag has a higher priority than configuration and environmnet variables.
	clientCmd.Flags().StringVarP(&serverLink, "link", "l", "", "CertArk server link")

	rootCmd.AddCommand(clientCmd)
}
