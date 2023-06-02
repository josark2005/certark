package certark

import (
	"encoding/json"
	"os"
)

func IsSystemd() bool {
	if _, err := os.Stat("/run/systemd/system"); err == nil {
		return true
	}
	return false
}

func FileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func ReadFileAndParseJson(path string, v any) error {
	// open and read file
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, v)
}
