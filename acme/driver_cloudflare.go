package acme

import (
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/tidwall/gjson"
)

func init() {
	Driver["cloudflare"] = &CloudflareDriver{}
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

func (c *CloudflareDriver) ReadConfFromJson(json string) *CloudflareDriver {
	c.Config.AuthEmail = gjson.Get(json, "dns_authuser").String()
	c.Config.AuthKey = gjson.Get(json, "dns_authkey").String()
	c.Config.AuthToken = gjson.Get(json, "dns_authtoken").String()
	c.Config.ZoneToken = gjson.Get(json, "dns_zonetoken").String()
	c.Config.TTL = int(gjson.Get(json, "dns_ttl").Int())
	c.Config.PropagationTimeout = time.Millisecond * time.Duration(gjson.Get(json, "dns_propagation_timeout").Int())
	c.Config.PollingInterval = time.Millisecond * time.Duration(gjson.Get(json, "dns_polling_interval").Int())
	return c
}
