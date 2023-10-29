package acme

import (
	"crypto"

	"github.com/go-acme/lego/v4/registration"
)

const (
	MODE_STAGING    = 0
	MODE_PRODUCTION = 1
)

type AcmeUser struct {
	Email        string
	Registration *registration.Resource
	Key          crypto.PrivateKey
}
