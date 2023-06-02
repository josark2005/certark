package cmd

import (
	"errors"
	"os"

	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
)

func init() {

	force_init_flag := false

	// init command
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize CertArk",
		Run: func(cmd *cobra.Command, args []string) {
			(&InitRunCondition{
				ForceInit: force_init_flag,
			}).Run()
		},
	}
	initCmd.Flags().BoolVarP(&force_init_flag, "force", "f", false, "force initialize")

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
	ShowLog   bool
	ForceInit bool
}

func (r *InitRunCondition) Run() (bool, error) {
	// check if directory is lock
	if r.ShowLog {
		ark.Info().Str("path", initLockFilePath).Msg("Checking lock file")
	}

	if r.ForceInit {
		ark.Warn().Msg("Force initializing")
	}

	if certark.FileOrDirExists(initLockFilePath) && !r.CheckMode && !r.ForceInit {
		err := errors.New("config directory is locked")
		ark.Error().Err(err).Msg("Init condition check failed")
		return false, err
	}

	// check serviceConfigDir
	if r.ShowLog {
		ark.Info().Str("path", serviceConfigDir).Msg("Checking config directory")
	}
	if !certark.FileOrDirExists(serviceConfigDir) {
		if !r.CheckMode {
			err := os.MkdirAll(serviceConfigDir, os.ModePerm)
			if err != nil {
				ark.Error().Err(err).Msg("Run condition init failed")
				return false, err
			}
		} else {
			return false, errors.New(serviceConfigDir + " not found")
		}
	}

	// check serviceConfigPath
	if r.ShowLog {
		ark.Info().Str("path", serviceConfigPath).Msg("Checking service config file")
	}
	if !certark.FileOrDirExists(serviceConfigPath) {
		if !r.CheckMode {
			file, err := os.OpenFile(serviceConfigPath, os.O_WRONLY|os.O_CREATE, os.ModeExclusive)
			if err != nil {
				ark.Error().Err(err).Msg("Run condition init failed")
				return false, err
			}
			file.WriteString("")
		} else {
			return false, errors.New(serviceConfigPath + " not found")
		}
	}

	// check stateDir
	if r.ShowLog {
		ark.Info().Str("path", stateDir).Msg("Checking state directory")
	}
	if !certark.FileOrDirExists(stateDir) {
		if !r.CheckMode {
			err := os.MkdirAll(stateDir, os.ModePerm)
			if err != nil {
				ark.Error().Err(err).Msg("Run condition init failed")
				return false, err
			}
		} else {
			return false, errors.New(stateDir + " not found")
		}
	}

	// check taskConfigDir
	if r.ShowLog {
		ark.Info().Str("path", taskConfigDir).Msg("Checking task directory")
	}
	if !certark.FileOrDirExists(taskConfigDir) {
		if !r.CheckMode {
			err := os.MkdirAll(taskConfigDir, os.ModePerm)
			if err != nil {
				ark.Error().Err(err).Msg("Run condition init failed")
				return false, err
			}
		} else {
			return false, errors.New(taskConfigDir + " not found")
		}
	}

	// check acmeUserDir
	if r.ShowLog {
		ark.Info().Str("path", acmeUserDir).Msg("Checking acme user directory")
	}
	if !certark.FileOrDirExists(acmeUserDir) {
		if !r.CheckMode {
			err := os.MkdirAll(acmeUserDir, os.ModePerm)
			if err != nil {
				ark.Error().Err(err).Msg("Run condition init failed")
				return false, err
			}
		} else {
			return false, errors.New(acmeUserDir + " not found")
		}
	}

	// write lock file
	if !r.CheckMode {
		ark.Info().Str("path", initLockFilePath).Msg("Writing lock file")
		fp, err := os.OpenFile(initLockFilePath, os.O_CREATE, os.ModeExclusive)
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
	if !certark.FileOrDirExists(initLockFilePath) {
		ark.Warn().Str("file", initLockFilePath).Msg("Lock file does not exist")
		return false, errors.New("lock file does not exist")
	}

	ark.Info().Str("file", initLockFilePath).Msg("Removing lock file")

	err := os.Remove(initLockFilePath)
	if err != nil {
		ark.Error().Err(err).Msg("Remove lock file failed")
		return false, err
	}

	ark.Info().Str("file", initLockFilePath).Msg("Lock file removed")
	return true, nil
}
