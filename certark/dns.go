package certark

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"

	"github.com/jokin1999/certark/ark"
)

type DNSProfile struct {
	Enabled bool `json:"enabled"`

	Provider     string `json:"provider"`
	Account      string `json:"account"`
	APIKey       string `json:"api_key"`
	DNSAPIToken  string `json:"dns_api_key"`
	ZoneAPIToken string `json:"zone_api_token"`
}

var DefaultDNSUserProfile = DNSProfile{
	Enabled:      true,
	Provider:     "",
	Account:      "",
	APIKey:       "",
	DNSAPIToken:  "",
	ZoneAPIToken: "",
}

func GetDNSProfileContent(dns string) ([]byte, error) {
	// check file exists
	profile := DNSUserDir + "/" + dns
	if !FileOrDirExists(profile) || !IsFile(profile) {
		err := errors.New("DNS user " + dns + " does not exist")
		return []byte{}, err
	}

	// read file
	profileContent, err := os.ReadFile(profile)
	if err != nil {
		return []byte{}, err
	}

	return profileContent, err
}

func GetDNSProfile(dns string) (DNSProfile, error) {
	// read file
	profileContent, err := GetDNSProfileContent(dns)
	if err != nil {
		return DNSProfile{}, err
	}

	dnsProfile := DNSProfile{}
	err = json.Unmarshal(profileContent, &dnsProfile)
	if err != nil {
		return DNSProfile{}, err
	}

	return dnsProfile, err
}

func GetDNSProfileJSONHuman(dns string) (string, error) {
	// read file
	profileContent, err := GetDNSProfileContent(dns)
	if err != nil {
		return "", err
	}

	var jsonBuff bytes.Buffer
	if err = json.Indent(&jsonBuff, profileContent, "", ""); err != nil {
		ark.Error().Err(err).Msg("Failed to show DNS user profile")
		return "", err
	}

	return jsonBuff.String(), nil
}
