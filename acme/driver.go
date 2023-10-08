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

// check acme driver exists
func IsDriverExists(driver string) bool {
	if _, ok := driverMap[driver]; ok {
		return true
	}
	return false
}
