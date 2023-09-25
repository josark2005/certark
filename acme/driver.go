package acme

import (
	"github.com/go-acme/lego/v4/challenge"
)

type ProviderDriver interface {
	NewDnsProviderConfig() (challenge.Provider, error)
	ReadConfFromJson(string)
}

type DriverConstruct func() ProviderDriver

var driverMap map[string]DriverConstruct = make(map[string]DriverConstruct)

func RegisterDriver(driverName string, driver DriverConstruct) {
	driverMap[driverName] = driver
}

//TODO -
