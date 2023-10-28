package cmd

import (
	"fmt"

	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
)

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
			listDnsUserProfiles()
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
				showDnsUserProfile(dns)
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
				addDnsUserProfile(dns)
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
		Use:   "set [DNS user profile]",
		Short: "Set config values in a DNS user profile",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) == 1 {
				dns := args[0]
				// set dns user profile enabled
				if cmd.Flags().Lookup("enable").Changed {
					setDnsUserProfile(dns, "enable", "true")
				}

				// set dns user profile disbaled
				if cmd.Flags().Lookup("disable").Changed {
					setDnsUserProfile(dns, "enable", "false")
				}

				// set dns user profile provider
				if cmd.Flags().Lookup("provider").Changed {
					setDnsUserProfile(dns, "provider", provider)
				}

				// set dns user profile account
				if cmd.Flags().Lookup("account").Changed {
					setDnsUserProfile(dns, "account", account)
				}

				// set dns user profile api_key
				if cmd.Flags().Lookup("apikey").Changed {
					setDnsUserProfile(dns, "api_key", api_key)
				}

				// set dns user profile dns_api_token
				if cmd.Flags().Lookup("dnstoken").Changed {
					setDnsUserProfile(dns, "dns_api_token", dns_api_token)
				}

				// set dns user profile zone_api_token
				if cmd.Flags().Lookup("zonetoken").Changed {
					setDnsUserProfile(dns, "zone_api_token", zone_api_token)
				}
			} else {
				cmd.Help()
			}
		},
	}

	c.Flags().BoolVar(&enabled, "enable", true, "enable DNS user profile")
	c.Flags().BoolVar(&enabled, "disable", false, "disable DNS user profile")

	c.Flags().StringVarP(&provider, "provider", "p", certark.DefaultDnsUserProfile.Provider, "set DNS provider")
	c.Flags().StringVarP(&account, "account", "a", certark.DefaultDnsUserProfile.Account, "set DNS provider account")
	c.Flags().StringVarP(&api_key, "apikey", "k", certark.DefaultDnsUserProfile.APIKey, "set DNS account API key")
	c.Flags().StringVarP(&dns_api_token, "dnstoken", "d", certark.DefaultDnsUserProfile.DNSAPIToken, "set DNS edit API token")
	c.Flags().StringVarP(&zone_api_token, "zonetoken", "z", certark.DefaultDnsUserProfile.ZoneAPIToken, "set DNS zone read API token")
	return c
}

// list dns user profiles
func listDnsUserProfiles() {
	dnsProfiles, err := certark.ListDnsUserProfiles()
	if err != nil {
		ark.Error().Err(err).Msg("Failed to list DNS profiles")
		return
	}
	for _, v := range dnsProfiles {
		fmt.Println(v)
	}
}

// show dns user profile
func showDnsUserProfile(dns string) {
	profile, err := certark.GetDnsJsonPretty(dns)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to show DNS user profile")
		return
	}
	fmt.Println(profile)
}

// add dns user profile
func addDnsUserProfile(dns string) {
	err := certark.AddDnsUser(dns)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to create dns user profile")
	}
}

// set dns user profile
func setDnsUserProfile(dns, key, value string) {
	err := certark.SetDnsUserProfile(dns, key, value)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to set dns user profile")
	}
}
