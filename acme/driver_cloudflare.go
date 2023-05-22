package acme

import (
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
)

func init() {
	driver := CloudflareDriver{}
	Driver["cloudflare"] = &driver
}

type CloudflareDriver struct {
	Config cloudflare.Config
}

func (c *CloudflareDriver) NewDnsProviderConfig() (challenge.Provider, error) {
	conf := cloudflare.Config{
		AuthEmail:          c.Config.AuthEmail,
		AuthKey:            c.Config.AuthKey,
		AuthToken:          c.Config.AuthToken,
		ZoneToken:          c.Config.ZoneToken,
		TTL:                c.Config.TTL,
		PropagationTimeout: c.Config.PropagationTimeout,
		PollingInterval:    c.Config.PollingInterval,
	}
	cf, err := cloudflare.NewDNSProviderConfig(&conf)
	return cf, err
}
