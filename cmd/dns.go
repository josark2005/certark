package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jokin1999/certark/acme"
	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
)

// check if dns profile exists
func checkDNSProfileExists(dns string) bool {
	res := certark.FileOrDirExists(certark.DNSUserDir + "/" + dns)
	if res {
		ark.Debug().Msg("DNS user profile exists")
	} else {
		ark.Debug().Msg("DNS user profile does not exist")
	}
	return res
}

func init() {
	// dns main command
	var dnsCmd = cmdDNS()

	// dns ls
	dnsCmd.AddCommand(cmdDNSLs())

	// dns show
	dnsCmd.AddCommand(cmdDNSShow())

	// dns add
	dnsCmd.AddCommand(cmdDNSAdd())

	// dns set
	dnsCmd.AddCommand(cmdDNSSet())

	rootCmd.AddCommand(dnsCmd)
}

// dns command
func cmdDNS() *cobra.Command {
	return &cobra.Command{
		Use:   "dns",
		Short: "DNS user profiles management",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
}

// dns ls
func cmdDNSLs() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List DNS user profiles",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			listDNSProfiles()
		},
	}
}

// dns show
func cmdDNSShow() *cobra.Command {
	return &cobra.Command{
		Use:   "show [TASK]",
		Short: "Show a task profile",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 0 {
				dns := args[0]
				showDNSUserProfile(dns)
			} else {
				cmd.Help()
			}
		},
	}
}

// dns add
func cmdDNSAdd() *cobra.Command {
	c := &cobra.Command{
		Use:   "add [DNS]",
		Short: "Add an DNS user profile",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 0 {
				dns := args[0]
				addDNSUserProfile(dns)
			} else {
				cmd.Help()
			}
		},
	}

	return c
}

// dns set
func cmdDNSSet() *cobra.Command {
	var (
		enabled        bool
		provider       string
		account        string
		api_key        string
		dns_api_token  string
		zone_api_token string
	)

	c := &cobra.Command{
		Use:   "set [DNS_user_profile]",
		Short: "Set config values in a DNS user profile",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) == 1 {
				dns := args[0]
				// set dns user profile enabled
				if cmd.Flags().Lookup("enable").Changed {
					ok := setDNSUserProfile(dns, "enable", "true")
					if !ok {
						ark.Error().Msg("Set domain failed")
					}
				}

				// set dns user profile disbaled
				if cmd.Flags().Lookup("disable").Changed {
					ok := setDNSUserProfile(dns, "enable", "false")
					if !ok {
						ark.Error().Msg("Set domain failed")
					}
				}

				// set dns user profile provider
				if cmd.Flags().Lookup("provider").Changed {
					ok := setDNSUserProfile(dns, "provider", provider)
					if !ok {
						ark.Error().Msg("Set provider failed")
					}
				}

				// set dns user profile account
				if cmd.Flags().Lookup("account").Changed {
					ok := setDNSUserProfile(dns, "account", account)
					if !ok {
						ark.Error().Msg("Set account failed")
					}
				}

				// set dns user profile api_key
				if cmd.Flags().Lookup("apikey").Changed {
					ok := setDNSUserProfile(dns, "api_key", api_key)
					if !ok {
						ark.Error().Msg("Set API key failed")
					}
				}

				// set dns user profile dns_api_token
				if cmd.Flags().Lookup("dnstoken").Changed {
					ok := setDNSUserProfile(dns, "dns_api_token", dns_api_token)
					if !ok {
						ark.Error().Msg("Set DNS edit API token failed")
					}
				}

				// set dns user profile zone_api_token
				if cmd.Flags().Lookup("zonetoken").Changed {
					ok := setDNSUserProfile(dns, "zone_api_token", zone_api_token)
					if !ok {
						ark.Error().Msg("Set DNS zone read API token failed")
					}
				}
			} else {
				cmd.Help()
			}
		},
	}

	c.Flags().BoolVar(&enabled, "enable", true, "enable task")
	c.Flags().BoolVar(&enabled, "disable", false, "disable task")

	c.Flags().StringVarP(&provider, "provider", "p", certark.DefaultDNSUserProfile.Provider, "set DNS provider")
	c.Flags().StringVarP(&account, "account", "a", certark.DefaultDNSUserProfile.Account, "set DNS provider account")
	c.Flags().StringVarP(&api_key, "apikey", "k", certark.DefaultDNSUserProfile.APIKey, "set DNS account API key")
	c.Flags().StringVarP(&dns_api_token, "dnstoken", "d", certark.DefaultDNSUserProfile.DNSAPIToken, "set DNS edit API token")
	c.Flags().StringVarP(&zone_api_token, "zonetoken", "z", certark.DefaultDNSUserProfile.ZoneAPIToken, "set DNS zone read API token")
	return c
}

// list dns user profiles
func listDNSProfiles() {
	err := filepath.Walk(certark.DNSUserDir, func(path string, info os.FileInfo, err error) error {
		if path == certark.TaskConfigDir {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		fmt.Println(path[len(certark.DNSUserDir)+1:])
		return nil
	})
	if err != nil {
		ark.Error().Err(err).Msg("Failed to list DNS user profiles")
		return
	}
}

// show dns user profile
func showDNSUserProfile(dns string) {
	profile, err := certark.GetDNSProfileJSONHuman(dns)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to show DNS user profile")
		return
	}
	fmt.Println(profile)
}

// add dns user profile
func addDNSUserProfile(dns string) {
	if checkDNSProfileExists(dns) {
		err := errors.New("DNS user existed")
		ark.Error().Err(err).Msg("Failed to create DNS user profile")
		return
	}

	// create profile
	fp, err := os.OpenFile(certark.DNSUserDir+"/"+dns, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to create DNS user profile")
		return
	}
	defer fp.Close()

	profile := certark.DefaultDNSUserProfile
	profileJson, _ := json.Marshal(profile)

	// write profile to file
	_, err = fp.WriteString(string(profileJson))
	if err != nil {
		ark.Error().Msg("Failed to create DNS user profile" + dns)
	} else {
		ark.Info().Msg("DNS user profile " + dns + " added")
	}
}

// set dns user profile
func setDNSUserProfile(dns, key, value string) bool {
	if !checkDNSProfileExists(dns) {
		err := errors.New("DNS user profile does not existed")
		ark.Error().Err(err).Msg("Failed to modify DNS user profile")
		return false
	}

	supportedKey := []string{
		"enable",
		"provider",
		"account",
		"api_key",
		"dns_api_key",
		"zone_api_token",
	}

	supportFlag := false // if option is supported, supportFlag will set to true
	for _, sk := range supportedKey {
		if key == sk {
			supportFlag = true
			break
		}
	}
	if !supportFlag {
		err := errors.New("not supported item")
		ark.Error().Str("key", key).Err(err).Msg("Failed to set DNS user profile")
		return false
	}

	ark.Info().Str("key", key).Str("value", value).Msg("Setting DNS user profile")

	// read profile
	profileContent, err := os.ReadFile(certark.DNSUserDir + "/" + dns)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to read DNS user profile")
	}
	ark.Debug().Str("content", string(profileContent)).Msg("Read DNS user profile")
	profile := certark.DNSProfile{}
	err = json.Unmarshal(profileContent, &profile)
	if err != nil {
		ark.Error().Err(err).Str("dns", dns).Msg("Failed to parse DNS user profile")
	}

	switch key {
	case "enable":
		if value == "true" {
			profile.Enabled = true
		} else {
			profile.Enabled = false
		}
	case "provider":
		if acme.IsDriverExists(value) {
			profile.Provider = value
		} else {
			ark.Error().Str("dns", dns).Msg("Failed to set provider")
			return false
		}
	case "account":
		profile.Account = value
	case "api_key":
		profile.APIKey = value
	case "dns_api_key":
		profile.DNSAPIToken = value
	case "zone_api_token":
		profile.ZoneAPIToken = value
	default:
		ark.Error().Msg("Failed to found a valid item")
		return false
	}

	// write profile to file
	profileJson, _ := json.Marshal(profile)
	fp, err := os.OpenFile(certark.DNSUserDir+"/"+dns, os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to open DNS user profile")
		return false
	}
	defer fp.Close()
	_, err = fp.WriteString(string(profileJson))
	if err != nil {
		ark.Error().Msg("Failed to modify DNS user profile " + dns)
	} else {
		ark.Info().Msg("DNS user profile " + dns + " modified")
	}

	return true
}
