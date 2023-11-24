package certark

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"regexp"
)

const reEmail = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`

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

func ReadFileAndParseJson(filepath string, v any) error {
	// open and read file
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, v)
}

func CheckEmail(email string) bool {
	exp, _ := regexp.Compile(reEmail)
	res := exp.Match([]byte(email))
	return res
}

func WriteStructToFile(s any, filepath string) error {
	profileJson, err := json.Marshal(s)
	if err != nil {
		return err
	}

	// create profile
	fp, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		err := errors.New("failed to create user profile")
		return err
	}
	defer fp.Close()

	// write profile to file
	_, err = fp.WriteString(string(profileJson))
	if err != nil {
		return err
	}
	return nil
}

// check if a json tag exists in a struct
func CheckStructJsonTagExists(s any, key string) bool {
	res := reflect.TypeOf(s)
	for i := 0; i < res.NumField(); i++ {
		if res.Field(i).Tag.Get("json") == key {
			return true
		}
	}
	return false
}
