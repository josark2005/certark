package acme

import (
	"log"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/jokin1999/certark/acme"
	"github.com/tidwall/gjson"
)

type CloudflareDriver struct {
	Config  cloudflare.Config
	Domains []string
	Bundle  bool // default: true
}

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

func (c *CloudflareDriver) ReadConfFromJson(json string) {
	c.Config.AuthEmail = gjson.Get(json, "dns_authuser").String()
	c.Config.AuthKey = gjson.Get(json, "dns_authkey").String()
	c.Config.AuthToken = gjson.Get(json, "dns_authtoken").String()
	c.Config.ZoneToken = gjson.Get(json, "dns_zonetoken").String()
	c.Config.TTL = int(gjson.Get(json, "dns_ttl").Int())
	c.Config.PropagationTimeout = time.Millisecond * time.Duration(gjson.Get(json, "dns_propagation_timeout").Int())
	c.Config.PollingInterval = time.Millisecond * time.Duration(gjson.Get(json, "dns_polling_interval").Int())

	// add domains
	for _, domain := range gjson.Get(json, "domains").Array() {
		c.Domains = append(c.Domains, domain.String())
	}

	// bundle certificate
	if gjson.Get(json, "cert_bundle").Exists() {
		c.Bundle = gjson.Get(json, "cert_bundle").Bool()
	} else {
		c.Bundle = true
	}
}

func (c *CloudflareDriver) Validate() (bool, error) {
	// API_Email, API_Key (dns_authmail, dns_authkey)
	// DNS_API_TOKEN (dns_authtoken)
	// DNS_API_TOKEN, ZONE_API_TOKEN (dns_authtoken, dns_zonetoken)
	// if (c.Config.AuthEmail != "" && c.Config.AuthKey != "") || c.Config.AuthToken != "" {
	// 	err := errors.New("invalid authentication method")
	// 	return false, err
	// }
	return true, nil
}

func (c *CloudflareDriver) RequestCertificate(client *lego.Client) (string, error) {
	request := certificate.ObtainRequest{
		Domains: c.Domains,
		Bundle:  c.Bundle,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}
	return string(certificates.Certificate), nil
}

// DON'T REMOVE, check implemention here
var _ acme.ProviderDriver = (*CloudflareDriver)(nil)
