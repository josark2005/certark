package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jokin1999/certark/acme"
	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/jokin1999/certark/tank"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
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
	return &cobra.Command{
		Use:   "acme",
		Short: "ACME configurations",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
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
			listAcmeUser()
		},
	}
}

// acme show
func cmdAcmeShow() *cobra.Command {
	return &cobra.Command{
		Use:   "show [EMAIL]",
		Short: "Show a acme user profile",
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
				acmeEmail := args[0]
				regAcmeUser(acmeEmail)
			} else {
				cmd.Help()
			}
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

// list acme users
func listAcmeUser() {
	err := filepath.Walk(acmeUserDir, func(path string, info os.FileInfo, err error) error {
		if path == acmeUserDir {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		fmt.Println(path[len(acmeUserDir)+1:])
		return nil
	})
	if err != nil {
		ark.Error().Err(err).Msg("Failed to list acme users")
		return
	}
}

// show acme user
func showAcmeUser(email string) {
	profile := acmeUserDir + "/" + email
	if !certark.FileOrDirExists(profile) || !certark.IsFile(profile) {
		err := errors.New("user " + email + " does not exist")
		ark.Error().Err(err).Msg("Failed to show acme user")
		return
	}

	// read file
	profileContent, err := os.ReadFile(profile)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to show acme user")
		return
	}

	var jsonBuff bytes.Buffer
	if err = json.Indent(&jsonBuff, profileContent, "", ""); err != nil {
		ark.Error().Err(err).Msg("Failed to show acme user")
		return
	}

	fmt.Println(jsonBuff.String())
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

// set acme user profile private key
func setAcmeUserPirvateKeyInFile(email string, privateKeyPath string) {
	if !checkUserExists(email) {
		// user does not exist
		err := errors.New("user does not exist")
		ark.Error().Err(err).Msg("Failed to set user profile")
		return
	}

	// set acme user profile
	fp, err := os.OpenFile(acmeUserDir+"/"+email, os.O_WRONLY|os.O_TRUNC, os.ModeExclusive)
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
		PrivateKey: string(bytes.Trim(privatekey, " \n")),
		Enabled:    true,
	}
	profileJson, _ := json.Marshal(profile)
	ark.Debug().Str("content", string(profileJson)).Msg("prepared profile data")

	// write profile to file
	_, err = fp.WriteString(string(profileJson))
	if err != nil {
		ark.Error().Err(err).Msg("Failed to add User " + email)
	} else {
		ark.Info().Msg("User " + email + " set")
	}
}

// register acme user
func regAcmeUser(email string) {
	if !checkUserExists(email) {
		// user does not exist
		err := errors.New("user does not exist")
		ark.Error().Err(err).Msg("Failed to find acme user profile")
		return
	}

	profilePath := acmeUserDir + "/" + email

	// read acme user profile
	profile, err := os.ReadFile(profilePath)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to read acme user profile")
		return
	}
	ark.Debug().Str("content", string(profile)).Msg("Read profile")

	// register acme user
	acmeUsername := email
	privateKey := ""
	if m, _ := tank.Load("MODE"); m == tank.MODE_PROD {
		privateKey = acme.RegisterAcmeUser(acmeUsername, acme.MODE_PRODUCTION)
	} else {
		privateKey = acme.RegisterAcmeUser(acmeUsername, acme.MODE_STAGING)
	}

	// regenerate profile
	newProfile := acmeUserProfile{
		Email:      gjson.Get(string(profile), "email").String(),
		PrivateKey: privateKey,
		Enabled:    gjson.Get(string(profile), "enabled").Bool(),
	}

	// write acme user profile
	fp, err := os.OpenFile(profilePath, os.O_WRONLY|os.O_TRUNC, os.ModeExclusive)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to open acme user profile")
		return
	}
	defer fp.Close()

	// convert json to string
	profileJson, _ := json.Marshal(newProfile)
	ark.Debug().Str("content", string(profileJson)).Msg("prepared profile data")

	// write profile
	_, err = fp.WriteString(string(profileJson))
	if err != nil {
		ark.Error().Err(err).Msg("Failed to write acme user profile")
		return
	}
}
