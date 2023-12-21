package cloudflare

import (
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/josark2005/certark/acme"
)

func init() {
	// register driver
	acme.RegisterDriver("cloudflare", func() acme.ProviderDriver {
		return &CloudflareDriver{}
	})
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

func (c *CloudflareDriver) LoadConf(domains []string, dns acme.DnsUserProfile) {
	c.Config.AuthEmail = dns.GetAccount()
	c.Config.AuthKey = dns.GetDnsApiToken()
	c.Config.AuthToken = dns.GetApiKey()
	c.Config.ZoneToken = dns.GetZoneApiToken()
	c.Config.TTL = int(dns.GetTTL())
	c.Config.PropagationTimeout = time.Second * time.Duration(dns.GetDnsPropagationTimeout())
	c.Config.PollingInterval = time.Second * time.Duration(dns.GetDnsPollingInterval())

	c.Domains = domains

	c.Bundle = dns.GetCertBundle()
}

func (c *CloudflareDriver) RequestCertificate(client *lego.Client) (string, error) {
	conf, err := c.NewDnsProviderConfig()
	if err != nil {
		return "", err
	}

	err = client.Challenge.SetDNS01Provider(conf)
	if err != nil {
		return "", err
	}

	request := certificate.ObtainRequest{
		Domains: c.Domains,
		Bundle:  c.Bundle,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return "", err
	}
	return string(certificates.Certificate), nil
}

// DON'T REMOVE, check implemention here
var _ acme.ProviderDriver = (*CloudflareDriver)(nil)
