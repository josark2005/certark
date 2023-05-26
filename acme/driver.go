package acme

import (
	"github.com/go-acme/lego/v4/challenge"
)

type DNSProvider interface {
	NewDnsProviderConfig() (challenge.Provider, error)
	ReadConfFromJson(string) *CloudflareDriver
}

var Driver map[string]DNSProvider = make(map[string]DNSProvider)
