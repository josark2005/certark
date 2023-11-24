package cloudflare

import "github.com/go-acme/lego/v4/providers/dns/cloudflare"

type CloudflareDriver struct {
	Config  cloudflare.Config
	Domains []string
	Bundle  bool // default: true
}
