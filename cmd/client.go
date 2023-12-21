package cmd

import (
	"github.com/josark2005/certark/ark"
	"github.com/josark2005/certark/certark"
	"github.com/spf13/cobra"
)

func init() {
	// CertArk server link
	var serverLink = ""
	var local = false

	var clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Start a CertArk client",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if cmd.Flags().Lookup("standalone").Changed && local {
				certark.Standalone(certark.ServiceConfigDir)
			}
		},
	}

	// Specify a server port
	clientCmd.Flags().StringVarP(&serverLink, "link", "l", "", "CertArk server link")

	// run in standalone mode
	clientCmd.Flags().BoolVarP(&local, "standalone", "a", true, "Run CertArk in standalone mode")

	rootCmd.AddCommand(clientCmd)
}
