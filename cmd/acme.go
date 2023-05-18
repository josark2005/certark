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
	var acmeEmail = ""
	c := &cobra.Command{
		Use:   "add",
		Short: "Add ACME user",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(acmeEmail) > 0 {
				if !checkEmail(acmeEmail) {
					ark.Warn().Msg("Unsupported email format")
				} else {
					// add acme user
					addAcmeUser(acmeEmail)
				}
			}
		},
	}
	c.Flags().StringVarP(&acmeEmail, "email", "e", "", "acme user email")
	c.MarkFlagRequired("email")
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
			}
		},
	}
	c.Flags().StringVarP(&acmeEmail, "email", "e", "", "acme user email")
	c.MarkFlagRequired("email")
	return c
}

// acme set
func cmdAcmeSet() *cobra.Command {
	var acmeEmail = ""
	c := &cobra.Command{
		Use:   "set [EMAIL] [PATH_OF_PRIVATE_KEY]",
		Short: "set ACME user",
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
	if checkUserExists(acmeUserDir + "/" + email) {
		// user exists
		err := errors.New("user existed")
		ark.Error().Err(err).Msg("Failed to create user profile")
		return
	}

	// create profile
	fp, err := os.OpenFile(acmeUserDir+"/"+email, os.O_CREATE|os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to create user profile")
		fp.Close()
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
	if !checkUserExists(acmeUserDir + "/" + email) {
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
	if !checkUserExists(acmeUserDir + "/" + email) {
		// user does not exist
		err := errors.New("user does not exist")
		ark.Error().Err(err).Msg("Failed to set user profile")
		return
	}

	// set acme user profile
	fp, err := os.OpenFile(acmeUserDir+"/"+email, os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to create user profile")
		fp.Close()
		return
	}
	defer fp.Close()

	// read private key
	// pkfp, err := os.OpenFile(privateKeyPath, os.O_RDONLY, os.ModeExclusive)
	// if err != nil {

	// }

	// prepare profile data
	profile := acmeUserProfile{
		Email:      email,
		PrivateKey: "",
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
