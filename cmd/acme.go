package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var acmeEmail = ""
	var acmeKeyFile = ""

	var acmeCmd = &cobra.Command{
		Use:   "acme",
		Short: "ACME configurations",
		Run: func(cmd *cobra.Command, args []string) {
			if acmeEmail == "" || acmeKeyFile == "" {
				cmd.Help()
			} else {
				fmt.Println(acmeEmail)
				fmt.Println(acmeKeyFile)
			}
		},
	}
	acmeCmd.Flags().StringVarP(&acmeEmail, "email", "e", "", "ACME user email")
	acmeCmd.Flags().StringVarP(&acmeKeyFile, "key", "k", "", "ACME user key")

	var acmeSetEmailCmd = &cobra.Command{
		Use:   "email",
		Short: "Set ACME user email",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				fmt.Println("Email:", args[0])
			} else {
				cmd.Help()
			}
		},
	}

	var acmeSetKeyCmd = &cobra.Command{
		Use:   "key",
		Short: "Set ACME user key file",
		Long:  "Specify a ACME user key file for CertArk to read, the key will be stored and managed by CertArk itself.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				fmt.Println("Key:", args[0])
			} else {
				cmd.Help()
			}
		},
	}

	acmeCmd.AddCommand(acmeSetEmailCmd)
	acmeCmd.AddCommand(acmeSetKeyCmd)

	rootCmd.AddCommand(acmeCmd)
}
