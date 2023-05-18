package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"

	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
)

const reEmail = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`

type acmeUserProfile struct {
	Email      string `json:"email"`
	PrivateKey string `json:"privatekey"`
	Enabled    bool   `json:"enabled"`
}

func checkEmail(email string) bool {
	exp, _ := regexp.Compile(reEmail)
	res := exp.Match([]byte(email))
	if res {
		ark.Debug().Msg("Email check passed")
	} else {
		ark.Debug().Msg("Email check failed")
	}
	return res
}

func init() {
	// acme main command
	var acmeCmd = cmdAcme()

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
	c := &cobra.Command{
		Use:   "add [EMAIL]",
		Short: "Add ACME user",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 0 {
				acmeEmail := args[0]
				if !checkEmail(acmeEmail) {
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
				rmAcmeUser(acmeEmail)
			}
		},
	}
	c.Flags().BoolVarP(&comfirm, "yes-i-really-mean-it", "", false, "comfirm to remove acme user")
	return c
}

// acme set
func cmdAcmeSet() *cobra.Command {
	//TODO - flags
	var acmeEmail = ""
	var acmePrivateKeyPath = ""
	c := &cobra.Command{
		Use:   "set [EMAIL] -k [PATH_OF_PRIVATE_KEY]",
		Short: "set ACME user",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(acmeEmail) > 0 && len(acmePrivateKeyPath) > 0 {
				setAcmeUser(acmeEmail, acmePrivateKeyPath)
			}
		},
	}
	c.Flags().StringVarP(&acmeEmail, "email", "e", "", "acme user email")
	c.Flags().StringVarP(&acmePrivateKeyPath, "key", "k", "", "file path of acme user private key")
	c.MarkFlagRequired("email")
	c.MarkFlagRequired("key")
	return c
}

// check if acme user exists
func checkUserExists(email string) bool {
	res := certark.FileOrDirExists(acmeUserDir + "/" + email)
	if res {
		ark.Debug().Msg("Acme user exists")
	} else {
		ark.Debug().Msg("Acme user does not exist")
	}
	return res
}

// add acme user
func addAcmeUser(email string) {
	if checkUserExists(email) {
		// user exists
		err := errors.New("user existed")
		ark.Error().Err(err).Msg("Failed to create user profile")
		return
	}

	// create profile
	fp, err := os.OpenFile(acmeUserDir+"/"+email, os.O_CREATE|os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to create user profile")
		return
	}
	defer fp.Close()

	profile := acmeUserProfile{
		Email:      email,
		PrivateKey: "",
		Enabled:    true,
	}
	profileJson, _ := json.Marshal(profile)

	// write profile to file
	_, err = fp.WriteString(string(profileJson))
	if err != nil {
		ark.Error().Msg("Failed to add User " + email)
	} else {
		ark.Info().Msg("User " + email + " added")
	}
}

// remove acme user
func rmAcmeUser(email string) {
	if !checkUserExists(email) {
		// user does not exist
		err := errors.New("user does not exist")
		ark.Error().Err(err).Msg("Failed to remove user profile")
		return
	}

	// remove profile
	err := os.Remove(acmeUserDir + "/" + email)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to remove user profile")
		return
	}

	ark.Info().Msg("User " + email + " removed")
}

// set acme user profile
func setAcmeUser(email string, privateKeyPath string) {
	if !checkUserExists(email) {
		// user does not exist
		err := errors.New("user does not exist")
		ark.Error().Err(err).Msg("Failed to set user profile")
		return
	}

	// set acme user profile
	fp, err := os.OpenFile(acmeUserDir+"/"+email, os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to set user profile")
		return
	}
	defer fp.Close()

	// read private key
	privatekey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to set user profile")
		return
	}

	// prepare profile data
	profile := acmeUserProfile{
		Email:      email,
		PrivateKey: string(privatekey),
		Enabled:    true,
	}
	profileJson, _ := json.Marshal(profile)

	// write profile to file
	_, err = fp.WriteString(string(profileJson))
	if err != nil {
		ark.Error().Err(err).Msg("Failed to add User " + email)
	} else {
		ark.Info().Msg("User " + email + " added")
	}
}
