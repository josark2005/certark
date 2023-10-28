package certark

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

// get acme user filepath
func GetTaskFilepath(name string) string {
	return TaskConfigDir + "/" + name
}

// check if acme user exists
func CheckTaskExists(name string) bool {
	res := FileOrDirExists(GetTaskFilepath(name))
	return res
}

// get task
func GetTask(name string) (TaskProfile, error) {
	profilePath := GetTaskFilepath(name)

	profile := TaskProfile{}
	err := ReadFileAndParseJson(profilePath, &profile)
	if err != nil {
		return TaskProfile{}, err
	}

	return profile, nil
}

// get task json
func GetTaskJson(name string) ([]byte, error) {
	profilePath := GetTaskFilepath(name)

	content, err := os.ReadFile(profilePath)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

// get task json pretty
func GetTaskJsonPretty(name string) (string, error) {
	profileContent, err := GetTaskJson(name)
	if err != nil {
		return "", err
	}

	var jsonBuff bytes.Buffer
	if err = json.Indent(&jsonBuff, profileContent, "", ""); err != nil {
		return "", err
	}

	return jsonBuff.String(), nil
}

// list tasks
func ListTasks() ([]string, error) {
	tasks := []string{}
	err := filepath.Walk(TaskConfigDir, func(path string, info os.FileInfo, err error) error {
		// skip dir itself
		if path == TaskConfigDir {
			return nil
		}
		// skip dirs
		if info.IsDir() {
			return nil
		}
		tasks = append(tasks, path[len(TaskConfigDir)+1:])
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return tasks, nil
}

// add acme user
func AddTask(name string) error {
	if CheckTaskExists(name) {
		err := errors.New("task exists")
		return err
	}

	profileFilepath := GetTaskFilepath(name)

	profile := DefaultTaskProfile

	err := WriteStructToFile(profile, profileFilepath)
	if err != nil {
		return err
	}
	return nil
}

// set task profile
func SetTaskProfile(name string, key string, value string) error {
	if !CheckTaskExists(name) {
		err := errors.New("task does not exist")
		return err
	}

	// check supported key
	if !CheckStructJsonTagExists(TaskProfile{}, key) {
		err := errors.New("task profile key not supported")
		return err
	}

	task, err := GetTask(name)
	if err != nil {
		return err
	}

	switch key {
	case "domain":
		task.Domain = []string{value}
	case "acme_user":
		if CheckAcmeUserExists(value) {
			task.AcmeUser = value
		} else {
			e := errors.New("failed to find acme user")
			return e
		}
	case "enable":
		if value == "true" {
			task.Enabled = true
		} else {
			task.Enabled = false
		}
	case "dns_profile":
		if CheckDnsUserExists(value) {
			task.DNSProfile = value
		} else {
			return errors.New("failed to find dns profile")
		}
		task.DNSProfile = value
	case "dns_ttl":
		v, e := strconv.Atoi(value)
		if e != nil {
			return e
		}
		task.DNSTTL = int64(v)
	case "dns_propagation_timeout":
		v, e := strconv.Atoi(value)
		if e != nil {
			return e
		}
		task.DNSPropagationTimeout = int64(v)
	case "dns_polling_interval":
		v, e := strconv.Atoi(value)
		if e != nil {
			return e
		}
		task.DNSPollingInterval = int64(v)
	case "url_check_enable":
		if value == "true" {
			task.UrlCheckEnable = true
		} else {
			task.UrlCheckEnable = false
		}
	case "url_check_target":
		task.UrlCheckTarget = value
	case "url_check_interval":
		v, e := strconv.Atoi(value)
		if e != nil {
			return e
		}
		task.UrlCheckInterval = int64(v)
	default:
		return errors.New("failed to found a valid item")
	}

	err = WriteStructToFile(task, GetTaskFilepath(name))
	if err != nil {
		return err
	}

	return nil
}

// append domains in a task profile
func AppendDomainTaskProfile(name string, domains []string) error {
	if !CheckTaskExists(name) {
		err := errors.New("task does not exist")
		return err
	}

	task, err := GetTask(name)
	if err != nil {
		return err
	}

	task.Domain = append(task.Domain, domains...)

	err = WriteStructToFile(task, GetTaskFilepath(name))
	if err != nil {
		return err
	}

	return nil
}

// remove domains in a task profile
func SubtractDomainTaskProfile(name string, domain string) error {
	if !CheckTaskExists(name) {
		err := errors.New("task does not exist")
		return err
	}

	task, err := GetTask(name)
	if err != nil {
		return err
	}

	domainsNew := []string{}
	for _, v := range task.Domain {
		if v != domain {
			domainsNew = append(domainsNew, v)
		} else {
			continue
		}
	}

	task.Domain = domainsNew

	err = WriteStructToFile(task, GetTaskFilepath(name))
	if err != nil {
		return err
	}
	return nil
}

// set acme user in a task profile
func SetAcmeUserTaskProfile(name string, acme string) error {
	if !CheckTaskExists(name) {
		err := errors.New("task does not exist")
		return err
	}

	// check if acme user exists
	if !CheckAcmeUserExists(acme) {
		err := errors.New("acme user does not existed")
		return err
	}

	task, err := GetTask(name)
	if err != nil {
		return err
	}

	task.AcmeUser = acme

	err = WriteStructToFile(task, GetTaskFilepath(name))
	if err != nil {
		return err
	}
	return nil
}
