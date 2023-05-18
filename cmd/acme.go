package cmd

import (
	"os"
	"regexp"

	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
)

func init() {
	// acme main command
	var acmeCmd = cmdAcme()

	acmeCmd.AddCommand(cmdAcmeAdd())
	acmeCmd.AddCommand(cmdAcmeRm())

	rootCmd.AddCommand(acmeCmd)
}

// acme command
func cmdAcme() *cobra.Command {
	return &cobra.Command{
		Use:   "acme",
		Short: "ACME configurations",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
}

// acme add command
func cmdAcmeAdd() *cobra.Command {
	var acmeEmail = ""
	c := &cobra.Command{
		Use:   "add",
		Short: "Add ACME user",
		Run: func(cmd *cobra.Command, args []string) {
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
	c.Flags().StringVarP(&acmeEmail, "email", "e", "", "acme user email")
	return c
}

// acme rm command
func cmdAcmeRm() *cobra.Command {
	var acmeEmail = ""
	c := &cobra.Command{
		Use:   "rm",
		Short: "Remove ACME user",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(acmeEmail) > 0 {
				rmAcmeUser(acmeEmail)
			} else {
				cmd.Help()
			}
		},
	}
	c.Flags().StringVarP(&acmeEmail, "email", "e", "", "acme user email")
	return c
}

// add acme user
func addAcmeUser(email string) {
	if certark.FileOrDirExists(acmeUserDir + "/" + email) {
		// user exists
		ark.Error().Str("error", "user existed").Msg("Failed to create user profile")
		return
	}

	// create profile
	fp, err := os.OpenFile(acmeUserDir+"/"+email, os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		ark.Error().Str("error", err.Error()).Msg("Failed to create user profile")
	}
	defer fp.Close()

	ark.Info().Msg("User " + email + " added")
}

// remove acme user
func rmAcmeUser(email string) {
	if !certark.FileOrDirExists(acmeUserDir + "/" + email) {
		// user does not exist
		ark.Error().Str("error", "user does not exist").Msg("Failed to remove user profile")
		return
	}

	// remove profile
	err := os.Remove(acmeUserDir + "/" + email)
	if err != nil {
		ark.Error().Str("error", err.Error()).Msg("Failed to remove user profile")
		return
	}

	ark.Info().Msg("User " + email + " removed")
}
