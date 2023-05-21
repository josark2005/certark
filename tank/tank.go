package tank

import (
	"errors"

	"github.com/jokin1999/certark/ark"
)

const (
	MODE_DEV  = "dev"
	MODE_PROD = "prod"
)

var Config = make(map[string]string)

func init() {
	Config["config_init"] = "true"
}

func Save(key, value string) {
	Config[key] = value
}

func Load(key string) (string, error) {
	value, ok := Config[key]
	if !ok {
		err := errors.New("key not found")
		ark.Error().Err(err).Str("key", key).Msg("Failed to load key")
		return "", err
	}
	return value, nil
}

func IsDevMode() bool {
	return Config["MODE"] == "dev"
}
