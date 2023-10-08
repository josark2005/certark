package certark

type DNSProfile struct {
	Enabled bool `json:"enabled"`

	Provider     string `json:"provider"`
	Account      string `json:"account"`
	APIKey       string `json:"api_key"`
	DNSAPIToken  string `json:"dns_api_key"`
	ZoneAPIToken string `json:"zone_api_token"`
}

var DefaultDNSUserProfile = DNSProfile{
	Enabled:      true,
	Provider:     "",
	Account:      "",
	APIKey:       "",
	DNSAPIToken:  "",
	ZoneAPIToken: "",
}
