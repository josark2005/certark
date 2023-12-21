package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/josark2005/certark/ark"
	"github.com/josark2005/certark/certark"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	// init command
	var initCmd = cmdInit()

	// init unlock command
	var initLockRemoveCmd = cmdInitUnlock()

	// add init unlock command to init command
	initCmd.AddCommand(initLockRemoveCmd)

	// add init cmd to root
	rootCmd.AddCommand(initCmd)

	// config command
	var configCmd = cmdConfig()

	// config show
	configCmd.AddCommand(cmdConfigShow())

	// config set
	configCmd.AddCommand(cmdConfigSet())

	// config current
	configCmd.AddCommand(cmdConfigCurrent())

	// add config command to root
	rootCmd.AddCommand(configCmd)
}

// init command
func cmdInit() *cobra.Command {
	force_init_flag := false
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize CertArk",
		Run: func(cmd *cobra.Command, args []string) {
			(&InitRunCondition{
				ForceInit: force_init_flag,
			}).Run()
		},
	}
	initCmd.Flags().BoolVarP(&force_init_flag, "force", "", false, "force initialize")
	return initCmd
}

// init unlock command
func cmdInitUnlock() *cobra.Command {
	var remove_confirm_flag = false
	var initLockRemoveCmd = &cobra.Command{
		Use:   "unlock",
		Short: "Remove init lock file",
		Run: func(cmd *cobra.Command, args []string) {
			// check comfirm flag
			if cmd.Flags().Lookup("yes-i-really-mean-it").Changed {
				// remove lock file
				removeLockfile()
			} else {
				ark.Warn().Msg("A comfirm flag is required, add --yes-i-really-mean-it flag at the end of the command")
			}
		},
	}
	initLockRemoveCmd.Flags().BoolVarP(&remove_confirm_flag, "yes-i-really-mean-it", "", false, "comfirm to remove init lock")
	return initLockRemoveCmd
}

type InitRunCondition struct {
	// In check mode, files will not be changed.
	CheckMode bool

	// show logs
	ShowLog bool

	// force initializing flag
	ForceInit bool
}

func (r *InitRunCondition) Run() (bool, error) {
	// check if directory is lock
	if r.ShowLog {
		ark.Info().Str("path", certark.InitLockFilePath).Msg("Checking lock file")
	}

	if r.ForceInit {
		ark.Warn().Msg("Force initializing")
	}

	if certark.FileOrDirExists(certark.InitLockFilePath) && !r.CheckMode && !r.ForceInit {
		err := errors.New("config directory is locked")
		ark.Error().Err(err).Msg("Init condition check failed")
		return false, err
	}

	// check serviceConfigDir
	if r.ShowLog {
		ark.Info().Str("path", certark.ServiceConfigDir).Msg("Checking config directory")
	}
	if !certark.FileOrDirExists(certark.ServiceConfigDir) {
		if !r.CheckMode {
			err := os.MkdirAll(certark.ServiceConfigDir, os.ModePerm)
			if err != nil {
				ark.Error().Err(err).Msg("Run condition init failed")
				return false, err
			}
		} else {
			return false, errors.New(certark.ServiceConfigDir + " not found")
		}
	}

	// check serviceConfigPath
	if r.ShowLog {
		ark.Info().Str("path", certark.ServiceConfigPath).Msg("Checking service config file")
	}
	if !certark.FileOrDirExists(certark.ServiceConfigPath) {
		if !r.CheckMode {
			file, err := os.OpenFile(certark.ServiceConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
			if err != nil {
				ark.Error().Err(err).Msg("Run condition init failed")
				return false, err
			}
			// fetch default config
			profileYaml, _ := yaml.Marshal(certark.DefaultConfig)
			file.WriteString(string(profileYaml))
		} else {
			return false, errors.New(certark.ServiceConfigPath + " not found")
		}
	}

	// check stateDir
	if r.ShowLog {
		ark.Info().Str("path", certark.StateDir).Msg("Checking state directory")
	}
	if !certark.FileOrDirExists(certark.StateDir) {
		if !r.CheckMode {
			err := os.MkdirAll(certark.StateDir, os.ModePerm)
			if err != nil {
				ark.Error().Err(err).Msg("Run condition init failed")
				return false, err
			}
		} else {
			return false, errors.New(certark.StateDir + " not found")
		}
	}

	// check taskConfigDir
	if r.ShowLog {
		ark.Info().Str("path", certark.TaskConfigDir).Msg("Checking task directory")
	}
	if !certark.FileOrDirExists(certark.TaskConfigDir) {
		if !r.CheckMode {
			err := os.MkdirAll(certark.TaskConfigDir, os.ModePerm)
			if err != nil {
				ark.Error().Err(err).Msg("Run condition init failed")
				return false, err
			}
		} else {
			return false, errors.New(certark.TaskConfigDir + " not found")
		}
	}

	// check acmeUserDir
	if r.ShowLog {
		ark.Info().Str("path", certark.AcmeUserDir).Msg("Checking acme user directory")
	}
	if !certark.FileOrDirExists(certark.AcmeUserDir) {
		if !r.CheckMode {
			err := os.MkdirAll(certark.AcmeUserDir, os.ModePerm)
			if err != nil {
				ark.Error().Err(err).Msg("Run condition init failed")
				return false, err
			}
		} else {
			return false, errors.New(certark.AcmeUserDir + " not found")
		}
	}

	// check DNSUserDir
	if r.ShowLog {
		ark.Info().Str("path", certark.DnsUserDir).Msg("Checking dns user directory")
	}
	if !certark.FileOrDirExists(certark.DnsUserDir) {
		if !r.CheckMode {
			err := os.MkdirAll(certark.DnsUserDir, os.ModePerm)
			if err != nil {
				ark.Error().Err(err).Msg("Run condition init failed")
				return false, err
			}
		} else {
			return false, errors.New(certark.DnsUserDir + " not found")
		}
	}

	// write lock file
	if !r.CheckMode {
		ark.Info().Str("path", certark.InitLockFilePath).Msg("Writing lock file")
		fp, err := os.OpenFile(certark.InitLockFilePath, os.O_CREATE|os.O_TRUNC, 0444)
		if err != nil {
			ark.Error().Err(err).Msg("Run condition init failed")
		}
		defer fp.Close()

		ark.Info().Msg("Initialization success")
	}

	return true, nil
}

func CheckRunCondition() bool {
	r, _ := (&InitRunCondition{CheckMode: true, ShowLog: false}).Run()
	return r
}

func CheckRunConditionWithLog() bool {
	r, _ := (&InitRunCondition{CheckMode: true, ShowLog: true}).Run()
	return r
}

func removeLockfile() (bool, error) {
	// check if lock file exists
	if !certark.FileOrDirExists(certark.InitLockFilePath) {
		ark.Warn().Str("file", certark.InitLockFilePath).Msg("Lock file does not exist")
		return false, errors.New("lock file does not exist")
	}

	ark.Info().Str("file", certark.InitLockFilePath).Msg("Removing lock file")

	err := os.Remove(certark.InitLockFilePath)
	if err != nil {
		ark.Error().Err(err).Msg("Remove lock file failed")
		return false, err
	}

	ark.Info().Str("file", certark.InitLockFilePath).Msg("Lock file removed")
	return true, nil
}

// config command
func cmdConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "CertArk settings",
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
	configFile := certark.ServiceConfigPath

	// read file
	profileContent, err := os.ReadFile(configFile)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to read config file")
		return
	}

	fmt.Println(string(profileContent))
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

	c.Flags().StringVarP(&mode, "mode", "m", "keep", "Set CertArk running mode: prod, dev")
	c.Flags().Int64VarP(&port, "port", "p", 0, "Set CertArk running port (server)")

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
	config, err := certark.ReadConfig(false)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to parse config file")
		return false
	}

	ark.Debug().Str("opt", option).Str("val", value).Msg("Setting config")

	// change config
	switch option {
	case "mode":
		if value == certark.MODE_DEV || value == certark.MODE_PROD {
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
	ark.Debug().Msg(string(profileYaml))
	fp, err := os.OpenFile(certark.ServiceConfigPath, os.O_WRONLY|os.O_TRUNC, 0660)
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

// config current
func cmdConfigCurrent() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show current CertArk configuration (current running)",
		Run: func(cmd *cobra.Command, args []string) {
			certark.LoadConfig(false)
			showCurrentConfig()
		},
	}
}

// show current config
func showCurrentConfig() {
	config, _ := yaml.Marshal(certark.CurrentConfig)
	fmt.Println(string(config))
}
