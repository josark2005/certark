package cmd

import (
	"os"
	"regexp"

	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
)

func init() {
	var acmeEmail = ""
	// var acmeKeyFile = ""

	var acmeCmd = &cobra.Command{
		Use:   "acme",
		Short: "ACME configurations",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var acmeAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add ACME user",
		Run: func(cmd *cobra.Command, args []string) {
			ark.Info().Msg("Adding user ")
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(acmeEmail) > 0 {
				exp, _ := regexp.Compile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
				if !exp.Match([]byte(acmeEmail)) {
					ark.Warn().Msg("Unsupported email format")
				} else {
					// add acme user
					addAcmeUser(acmeEmail)
				}
			} else {
				cmd.Help()
			}
		},
	}
	acmeAddCmd.Flags().StringVarP(&acmeEmail, "email", "e", "", "acme user email")

	// var acmeSetKeyCmd = &cobra.Command{
	// 	Use:   "key",
	// 	Short: "Set ACME user key file",
	// 	Long:  "Specify a ACME user key file for CertArk to read, the key will be stored and managed by CertArk itself.",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		if len(args) > 0 {
	// 			fmt.Println("Key:", args[0])
	// 		} else {
	// 			cmd.Help()
	// 		}
	// 	},
	// }

	acmeCmd.AddCommand(acmeAddCmd)
	// acmeCmd.AddCommand(acmeSetKeyCmd)

	rootCmd.AddCommand(acmeCmd)
}

func addAcmeUser(email string) {
	if certark.FileOrDirExists(acmeUserDir + "/" + email) {
		// user exists
	}

	// create profile
	fp, err := os.OpenFile(acmeUserDir+"/"+email, os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		ark.Error().Str("error", err.Error()).Msg("Failed to create config file")
	}
	defer fp.Close()
}
