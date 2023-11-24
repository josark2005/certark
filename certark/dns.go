package certark

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jokin1999/certark/acme"
)

// get dns profile filepath
func GetDnsFilepath(name string) string {
	return DnsUserDir + "/" + name
}

// check if dns profile exists
func CheckDnsUserExists(name string) bool {
	res := FileOrDirExists(GetDnsFilepath(name))
	return res
}

// get dns profile
func GetDns(name string) (DnsUserProfile, error) {
	if !CheckDnsUserExists(name) {
		err := errors.New("dns profile does not exist")
		return DnsUserProfile{}, err
	}

	profilePath := GetDnsFilepath(name)

	profile := DnsUserProfile{}
	err := ReadFileAndParseJson(profilePath, &profile)
	if err != nil {
		return DnsUserProfile{}, err
	}

	return profile, nil
}

// get dns profile json
func GetDnsJson(name string) ([]byte, error) {
	if !CheckDnsUserExists(name) {
		err := errors.New("dns profile does not exist")
		return []byte{}, err
	}

	profilePath := GetDnsFilepath(name)

	content, err := os.ReadFile(profilePath)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

// get dns profile json pretty
func GetDnsJsonPretty(name string) (string, error) {
	if !CheckDnsUserExists(name) {
		err := errors.New("dns profile does not exist")
		return "", err
	}

	profileContent, err := GetDnsJson(name)
	if err != nil {
		return "", err
	}

	var jsonBuff bytes.Buffer
	if err = json.Indent(&jsonBuff, profileContent, "", ""); err != nil {
		return "", err
	}

	return jsonBuff.String(), nil
}

// list dns user profiles
func ListDnsUserProfiles() ([]string, error) {
	dnsProfiles := []string{}
	err := filepath.Walk(DnsUserDir, func(path string, info os.FileInfo, err error) error {
		// skip dir itself
		if path == DnsUserDir {
			return nil
		}
		// skip dirs
		if info.IsDir() {
			return nil
		}
		dnsProfiles = append(dnsProfiles, path[len(DnsUserDir)+1:])
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return dnsProfiles, nil
}

// add dns user profile
func AddDnsUser(name string) error {
	if CheckDnsUserExists(name) {
		err := errors.New("dns user exists")
		return err
	}

	profileFilepath := GetDnsFilepath(name)

	profile := DefaultDnsUserProfile

	err := WriteStructToFile(profile, profileFilepath)
	if err != nil {
		return err
	}
	return nil
}

// set dns user profile
func SetDnsUserProfile(name string, key string, value string) error {
	if !CheckDnsUserExists(name) {
		err := errors.New("dns user does not exist")
		return err
	}

	// check supported key
	if !CheckStructJsonTagExists(DnsUserProfile{}, key) {
		err := errors.New("dns user profile key not supported")
		return err
	}

	dns, err := GetDns(name)
	if err != nil {
		return err
	}

	switch key {
	case "enabled":
		if value == "true" {
			dns.Enabled = true
		} else {
			dns.Enabled = false
		}
	case "provider":
		if acme.CheckDriverExists(value) {
			dns.Provider = value
		}
	case "account":
		dns.Account = value
	case "api_key":
		dns.ApiKey = value
	case "dns_api_token":
		dns.ApiKey = value
	case "zone_api_token":
		dns.ZoneApiToken = value
	case "dns_ttl":
		v, e := strconv.Atoi(value)
		if e != nil {
			return e
		}
		dns.DnsTTL = int64(v)
	case "dns_propagation_timeout":
		v, e := strconv.Atoi(value)
		if e != nil {
			return e
		}
		dns.DnsPropagationTimeout = int64(v)
	case "dns_polling_interval":
		v, e := strconv.Atoi(value)
		if e != nil {
			return e
		}
		dns.DnsPollingInterval = int64(v)
	default:
		return errors.New("failed to found a valid item")
	}

	err = WriteStructToFile(dns, GetDnsFilepath(name))
	if err != nil {
		return err
	}

	return nil
}
