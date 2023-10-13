package acme

import (
	"errors"
	"fmt"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/lego"
	// _ "github.com/jokin1999/certark/acme/drivers"
)

type ProviderDriver interface {
	NewDnsProviderConfig() (challenge.Provider, error)
	ReadConfFromJson(string)
	Validate() (bool, error)
	RequestCertificate(*lego.Client) (string, error)
}

type DriverConstructor func() ProviderDriver

// drivers list
var driverMap = map[string]DriverConstructor{}

// reigster driver
func RegisterDriver(driverName string, driver DriverConstructor) {
	driverMap[driverName] = driver
	fmt.Println(driverMap)
}

// check acme driver exists
func IsDriverExists(driver string) bool {
	_, ok := driverMap[driver]
	return ok
}

func GetDriver(driverName string) (DriverConstructor, error) {
	if !IsDriverExists(driverName) {
		err := errors.New("driver not found: " + driverName)
		return nil, err
	}
	return driverMap[driverName], nil
}

// import drivers
// var _ := drivers
