package acme

import (
	"errors"
	// _ "github.com/jokin1999/certark/acme/drivers"
)

// reigster driver
func RegisterDriver(driverName string, driver DriverConstructor) {
	driverMap[driverName] = driver
}

// check acme driver exists
func CheckDriverExists(driver string) bool {
	_, ok := driverMap[driver]
	return ok
}

func GetDriver(driverName string) (DriverConstructor, error) {
	if driverName == "" {
		err := errors.New("empty driver name is not supported")
		return nil, err
	}
	if !CheckDriverExists(driverName) {
		err := errors.New("driver not found: " + driverName)
		return nil, err
	}
	return driverMap[driverName], nil
}
