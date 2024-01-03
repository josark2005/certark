package cmd

import (
	"fmt"

	"github.com/josark2005/certark/ark"
	"github.com/josark2005/certark/certark"
	"github.com/spf13/cobra"
)

func init() {
	// acme main command
	var acmeCmd = cmdAcme()

	// acme ls
	acmeCmd.AddCommand(cmdAcmeLs())

	// acme show
	acmeCmd.AddCommand(cmdAcmeShow())

	// acme reg
	acmeCmd.AddCommand(cmdAcmeReg())

	// acme add
	acmeCmd.AddCommand(cmdAcmeAdd())

	// acme rm
	acmeCmd.AddCommand(cmdAcmeRm())

	// acme set (profile)
	acmeCmd.AddCommand(cmdAcmeSet())

	rootCmd.AddCommand(acmeCmd)
}

// acme command
func cmdAcme() *cobra.Command {
	var acme = &cobra.Command{
		Use:   "acme",
		Short: "ACME profiles management",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	return acme
}

// acme ls
func cmdAcmeLs() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List acme user profiles",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			listAcmeUsers()
		},
	}
}

// acme show
func cmdAcmeShow() *cobra.Command {
	return &cobra.Command{
		Use:     "show [EMAIL]",
		Short:   "Show a acme user profile",
		Aliases: []string{"inspec"},
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 0 {
				acmeEmail := args[0]
				showAcmeUser(acmeEmail)
			} else {
				cmd.Help()
			}
		},
	}
}

// acme reg
func cmdAcmeReg() *cobra.Command {
	return &cobra.Command{
		Use:   "reg [EMAIL]",
		Short: "Register a acme user",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 0 {
				acme := args[0]
				regAcmeUser(acme)
			} else {
				cmd.Help()
			}
		},
	}
}

// acme add command
func cmdAcmeAdd() *cobra.Command {
	c := &cobra.Command{
		Use:   "add [ACME_USER_NAME] [EMAIL]",
		Short: "Add an ACME user",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) == 2 {
				acmeName := args[0]
				acmeEmail := args[1]
				if !certark.CheckEmail(acmeEmail) {
					ark.Warn().Msg("Unsupported email format")
				} else {
					// add acme user
					addAcmeUser(acmeName, acmeEmail)
				}
			} else {
				cmd.Help()
			}
		},
	}
	return c
}

// acme rm command
func cmdAcmeRm() *cobra.Command {
	var comfirm = false
	c := &cobra.Command{
		Use:   "rm [EMAIL]",
		Short: "Remove ACME user",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 0 {
				acmeEmail := args[0]
				if !comfirm {
					ark.Warn().Msg("A comfirm flag is required, add --yes-i-really-mean-it flag at the end of the command")
					return
				}
				removeAcmeUser(acmeEmail)
			}
		},
	}
	c.Flags().BoolVarP(&comfirm, "yes-i-really-mean-it", "", false, "comfirm to remove acme user")
	return c
}

// acme set
func cmdAcmeSet() *cobra.Command {
	var acmePrivateKeyPath = ""
	c := &cobra.Command{
		Use:   "set [EMAIL]",
		Short: "set ACME user",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 0 {
				if len(acmePrivateKeyPath) > 0 {
					acmeEmail := args[0]
					setAcmeUserPirvateKeyInFile(acmeEmail, acmePrivateKeyPath)
				}
			} else {
				cmd.Help()
			}
		},
	}
	c.Flags().StringVarP(&acmePrivateKeyPath, "key", "k", "", "file path of acme user private key")
	return c
}

// list acme users
func listAcmeUsers() {
	users, err := certark.ListAcmeUsers()
	if err != nil {
		ark.Error().Err(err).Msg("Failed to list acme users")
		return
	}
	for _, v := range users {
		fmt.Println(v)
	}
}

// show acme user
func showAcmeUser(acmeName string) {
	profile, err := certark.GetAcmeUserJsonPretty(acmeName)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to show acme user")
		return
	}
	fmt.Println(profile)
}

// add acme user
func addAcmeUser(name string, email string) {
	err := certark.AddAcmeUser(name, email)
	if err != nil {
		ark.Error().Msg("Failed to add User " + name)
	} else {
		ark.Info().Msg("User " + name + " added")
	}
}

// remove acme user
func removeAcmeUser(name string) {
	err := certark.RemoveAcmeUser(name)
	if err != nil {
		ark.Error().Msg("Failed to remove acme user " + name)
	} else {
		ark.Info().Msg("Acme user " + name + " removed")
	}

}

// set acme user profile private key
func setAcmeUserPirvateKeyInFile(acmeName string, privateKeyPath string) {
	err := certark.SetAcmeUserPrivateKeyInFile(acmeName, privateKeyPath)

	if err != nil {
		ark.Error().Err(err).Msg("Failed to update User " + acmeName)
	} else {
		ark.Info().Msg("User " + acmeName + " private key updated")
	}
}

// register acme user
func regAcmeUser(acmeName string) {
	err := certark.RegisterAcmeUser(acmeName)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to register acme user")
		return
	}
}
