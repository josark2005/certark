package cmd

import (
	"errors"
	"os"

	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
)

func init() {

	// init command
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize CertArk",
		Run: func(cmd *cobra.Command, args []string) {
			(&InitRunCondition{}).Run()
		},
	}

	var remove_confirm_flag = false
	var initLockRemoveCmd = &cobra.Command{
		Use:   "unlock",
		Short: "Remove init lock file",
		Run: func(cmd *cobra.Command, args []string) {
			// check comfirm flag
			if remove_confirm_flag {
				// remove lock file
				removeLockfile()
			} else {
				ark.Warn().Msg("A comfirm flag is required, add --yes-i-really-mean-it flag at the end of the command")
			}
		},
	}
	initLockRemoveCmd.Flags().BoolVarP(&remove_confirm_flag, "yes-i-really-mean-it", "", false, "comfirm to remove init lock")

	initCmd.AddCommand(initLockRemoveCmd)

	// service command
	//TODO - add service install

	rootCmd.AddCommand(initCmd)
}

type InitRunCondition struct {
	CheckMode bool
}

func (r *InitRunCondition) Run() (bool, error) {
	// check if directory is lock
	ark.Info().Str("path", initLockFilePath).Msg("Checking lock file")
	if certark.FileOrDirExists(initLockFilePath) && !r.CheckMode {
		ark.Error().Str("error", "config directory is locked").Msg("Init condition check failed")
		return false, errors.New("config directory is locked")
	}

	ark.Info().Str("path", serviceConfigDir).Msg("Checking config directory")
	if !certark.FileOrDirExists(serviceConfigDir) {
		if !r.CheckMode {
			err := os.MkdirAll(serviceConfigDir, os.ModePerm)
			if err != nil {
				ark.Error().Str("error", err.Error()).Msg("Run condition init failed")
				return false, err
			}
		} else {
			return false, errors.New(serviceConfigDir + " not found")
		}
	}

	ark.Info().Str("path", serviceConfigPath).Msg("Checking service config file")
	if !certark.FileOrDirExists(serviceConfigPath) {
		if !r.CheckMode {
			file, err := os.OpenFile(serviceConfigPath, os.O_WRONLY|os.O_CREATE, 0660)
			if err != nil {
				ark.Error().Str("error", err.Error()).Msg("Run condition init failed")
				return false, err
			}
			file.WriteString("")
		} else {
			return false, errors.New(serviceConfigPath + " not found")
		}
	}

	ark.Info().Str("path", domainConfigDir).Msg("Checking domain directory")
	if !certark.FileOrDirExists(domainConfigDir) {
		if !r.CheckMode {
			err := os.MkdirAll(domainConfigDir, os.ModePerm)
			if err != nil {
				ark.Error().Str("error", err.Error()).Msg("Run condition init failed")
				return false, err
			}
		} else {
			return false, errors.New(domainConfigDir + " not found")
		}
	}

	ark.Info().Str("path", acmeUserDir).Msg("Checking acme user directory")
	if !certark.FileOrDirExists(acmeUserDir) {
		if !r.CheckMode {
			err := os.MkdirAll(acmeUserDir, os.ModePerm)
			if err != nil {
				ark.Error().Str("error", err.Error()).Msg("Run condition init failed")
				return false, err
			}
		} else {
			return false, errors.New(acmeUserDir + " not found")
		}
	}

	// write lock file
	if !r.CheckMode {
		ark.Info().Str("path", initLockFilePath).Msg("Writing lock file")
		fp, err := os.OpenFile(initLockFilePath, os.O_CREATE, 0660)
		if err != nil {
			ark.Error().Str("error", err.Error()).Msg("Run condition init failed")
		}
		defer fp.Close()

		ark.Info().Msg("Initialization success")
	}

	return true, nil
}

func CheckRunCondition() bool {
	r, _ := (&InitRunCondition{CheckMode: true}).Run()
	return r
}

func removeLockfile() (bool, error) {
	// check if lock file exists
	if !certark.FileOrDirExists(initLockFilePath) {
		ark.Warn().Str("file", initLockFilePath).Msg("Lock file does not exist")
		return false, errors.New("lock file does not exist")
	}

	ark.Info().Str("file", initLockFilePath).Msg("Removing lock file")

	err := os.Remove(initLockFilePath)
	if err != nil {
		ark.Error().Str("error", err.Error()).Msg("Remove lock file failed")
		return false, err
	}

	ark.Info().Str("file", initLockFilePath).Msg("Lock file removed")
	return true, nil
}
