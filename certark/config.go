package certark

import (
	"os"

	"github.com/jokin1999/certark/ark"
	"gopkg.in/yaml.v3"
)

const (
	MODE_DEV  = "dev"
	MODE_PROD = "prod"
)

var (
	ServiceConfigDir  = "/etc/certark"
	InitLockFile      = ".lock"
	InitLockFilePath  = ServiceConfigDir + "/" + InitLockFile
	ServiceConfigFile = "config.yml"
	ServiceConfigPath = ServiceConfigDir + "/" + ServiceConfigFile
	TaskConfigDir     = ServiceConfigDir + "/task"
	StateDir          = ServiceConfigDir + "/state"
	AcmeUserDir       = ServiceConfigDir + "/user"
	CertarkService    = "certark.service"
)

type Config struct {
	Mode string `yaml:"mode"`
	Port int64  `yaml:"port"`
}

var DefaultConfig = Config{
	Mode: MODE_DEV,
	Port: 7701,
}

var CurrentConfig = Config{}

// read config
func ReadConfig(slient bool) (Config, error) {
	configFile := ServiceConfigPath

	// read file
	profileContent, err := os.ReadFile(configFile)
	if err != nil {
		if !slient {
			ark.Warn().Err(err).Msg("Failed to read config file")
		}
		return Config{}, err
	}

	// parse
	config := Config{}
	err = yaml.Unmarshal(profileContent, &config)
	return config, err
}

// load config
func LoadConfig(slient bool) {
	config, err := ReadConfig(slient)
	if err != nil {
		if !slient {
			ark.Warn().Err(err).Msg("Load CertArk config failed, may fallback to default")
		}
	}
	CurrentConfig = config
}
