package cmd

import (
	"errors"
	"os"

	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
)

const (
	serviceConfigDir  = "/etc/certark"
	serviceConfigFile = "config.yml"
	serviceConfigPath = serviceConfigDir + "/" + serviceConfigFile
	domainConfigDir   = serviceConfigDir + "/domain"
	acmeUserKeyDir    = serviceConfigDir + "/userkey"
	certarkService    = "certark.service"
)

func init() {

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize CertArk",
		Run: func(cmd *cobra.Command, args []string) {
			(&InitRunCondition{}).Run()
		},
	}

	rootCmd.AddCommand(initCmd)
}

type InitRunCondition struct {
	CheckMode bool
}

func (r *InitRunCondition) Run() (bool, error) {
	if !certark.FileOrDirExists(serviceConfigDir) {
		if !r.CheckMode {
			err := os.MkdirAll(serviceConfigDir, os.ModePerm)
			if err != nil {
				ark.Warn().Str("error", err.Error()).Msg("Run condition check failed")
				return false, err
			}
		} else {
			return false, errors.New(serviceConfigDir + " not found")
		}
	}

	if !certark.FileOrDirExists(serviceConfigPath) {
		if !r.CheckMode {
			file, err := os.OpenFile(serviceConfigPath, os.O_WRONLY|os.O_CREATE, 0660)
			if err != nil {
				ark.Warn().Str("error", err.Error()).Msg("Run condition check failed")
				return false, err
			}
			file.WriteString("")
		} else {
			return false, errors.New(serviceConfigPath + " not found")
		}
	}

	if !certark.FileOrDirExists(domainConfigDir) {
		if !r.CheckMode {
			err := os.MkdirAll(domainConfigDir, os.ModePerm)
			if err != nil {
				ark.Warn().Str("error", err.Error()).Msg("Run condition check failed")
				return false, err
			}
		} else {
			return false, errors.New(domainConfigDir + " not found")
		}
	}

	if !certark.FileOrDirExists(acmeUserKeyDir) {
		if !r.CheckMode {
			err := os.MkdirAll(acmeUserKeyDir, os.ModePerm)
			if err != nil {
				ark.Warn().Str("error", err.Error()).Msg("Run condition check failed")
				return false, err
			}
		} else {
			return false, errors.New(acmeUserKeyDir + " not found")
		}
	}

	return true, nil
}

func CheckRunCondition() bool {
	r, _ := (&InitRunCondition{CheckMode: true}).Run()
	return r
}
