package acme

import (
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/lego"
)

type DnsUserProfile interface {
	GetProvider() string
	GetAccount() string
	GetApiKey() string
	GetDnsApiToken() string
	GetZoneApiToken() string
	GetTTL() int64
	GetDnsPropagationTimeout() int64
	GetDnsPollingInterval() int64
	GetCertBundle() bool
}

type ProviderDriver interface {
	NewDnsProviderConfig() (challenge.Provider, error)
	LoadConf([]string, DnsUserProfile)
	RequestCertificate(*lego.Client) (string, error)
}

type DriverConstructor func() ProviderDriver

// drivers list
var driverMap = map[string]DriverConstructor{}
