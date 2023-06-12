package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	// config main command
	var configCmd = cmdConfig()

	// config show
	configCmd.AddCommand(cmdConfigShow())

	// config set
	configCmd.AddCommand(cmdConfigSet())

	// config current
	configCmd.AddCommand(cmdConfigCurrent())

	rootCmd.AddCommand(configCmd)
}

// config command
func cmdConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "CertArk configurations",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
}

// config show
func cmdConfigShow() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show CertArk configuration",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			showConfig()
		},
	}
}

// show config
func showConfig() {
	configFile := serviceConfigPath

	// read file
	profileContent, err := os.ReadFile(configFile)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to read config file")
		return
	}

	fmt.Println(string(profileContent))
}

// read config
func ReadConfig() (certark.Config, error) {
	configFile := serviceConfigPath

	// read file
	profileContent, err := os.ReadFile(configFile)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to read config file")
		return certark.Config{}, err
	}

	// parse
	config := certark.Config{}
	err = yaml.Unmarshal(profileContent, &config)
	return config, err
}

// config set
func cmdConfigSet() *cobra.Command {
	var mode string
	var port int64

	c := &cobra.Command{
		Use:   "set",
		Short: "Set CertArk configuration",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			// mode
			if cmd.Flags().Lookup("mode").Changed && mode != "keep" {
				setConfig("mode", mode)
			}

			// port
			if cmd.Flags().Lookup("port").Changed && port != 0 {
				setConfig("port", strconv.Itoa(int(port)))
			}
		},
	}

	c.Flags().StringVarP(&mode, "mode", "m", "keep", "Set CertArk running mode")
	c.Flags().Int64VarP(&port, "port", "p", 0, "Set CertArk running mode")

	return c
}

// set config
func setConfig(option string, value string) bool {
	var supportFlag = false
	var supportedOption = []string{
		"mode",
		"port",
	}

	for _, v := range supportedOption {
		if v == option {
			supportFlag = true
			break
		}
	}

	if !supportFlag {
		ark.Warn().Str("opt", option).Msg("Unsupported config")
		return false
	}

	// read config
	configFile := serviceConfigPath
	profileContent, err := os.ReadFile(configFile)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to read config file")
		return false
	}
	config := certark.Config{}
	err = yaml.Unmarshal(profileContent, &config)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to parse config file")
		return false
	}

	ark.Debug().Str("opt", option).Str("val", value).Msg("Setting config")

	// change config
	switch option {
	case "mode":
		if value == "dev" || value == "prod" {
			config.Mode = value
		} else {
			ark.Warn().Str("opt", option).Str("val", value).Msg("Unsupported config value")
			return false
		}
	case "port":
		v, e := strconv.Atoi(value)
		if e != nil {
			ark.Error().Err(err).Msg("Invalid port number")
			return false
		}
		config.Port = int64(v)
	}

	// write back
	profileYaml, _ := yaml.Marshal(config)
	fmt.Println(string(profileYaml))
	fp, err := os.OpenFile(configFile, os.O_WRONLY|os.O_TRUNC, os.ModeExclusive)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to open config file")
		return false
	}
	defer fp.Close()
	_, err = fp.WriteString(string(profileYaml))
	if err != nil {
		ark.Error().Msg("Failed to save config")
	} else {
		ark.Info().Msg("Save config success")
	}

	return true
}

// config cu
func cmdConfigCurrent() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show current CertArk configuration (current running)",
		Run: func(cmd *cobra.Command, args []string) {
			showCurrentConfig()
		},
	}
}

// show current config
func showCurrentConfig() {
	config, _ := yaml.Marshal(certark.CurrentConfig)
	fmt.Println(string(config))
}
