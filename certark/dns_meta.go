package certark

type DnsUserProfile struct {
	Enabled bool `json:"enabled"`

	Provider              string `json:"provider"`
	Account               string `json:"account"`
	ApiKey                string `json:"api_key"`
	DnsApiToken           string `json:"dns_api_token"`
	ZoneApiToken          string `json:"zone_api_token"`
	DnsTTL                int64  `json:"dns_ttl"`                 // ttl 120 is recommanded
	DnsPropagationTimeout int64  `json:"dns_propagation_timeout"` // in millisecond, 60*1000 is recommanded
	DnsPollingInterval    int64  `json:"dns_polling_interval"`    // in millisecond, 5 *1000 is recommanded
	CertBundle            bool   `json:"cert_bundle"`
}

var DefaultDnsUserProfile = DnsUserProfile{
	Enabled:               true,
	Provider:              "",
	Account:               "",
	ApiKey:                "",
	DnsApiToken:           "",
	ZoneApiToken:          "",
	DnsTTL:                120,
	DnsPropagationTimeout: 60,
	DnsPollingInterval:    5,
	CertBundle:            true,
}

func (d *DnsUserProfile) GetProvider() string {
	return d.Provider
}

func (d *DnsUserProfile) GetAccount() string {
	return d.Account
}

func (d *DnsUserProfile) GetApiKey() string {
	return d.ApiKey
}

func (d *DnsUserProfile) GetDnsApiToken() string {
	return d.DnsApiToken
}

func (d *DnsUserProfile) GetZoneApiToken() string {
	return d.ZoneApiToken
}

func (d *DnsUserProfile) GetTTL() int64 {
	return d.DnsTTL
}

func (d *DnsUserProfile) GetDnsPropagationTimeout() int64 {
	return d.DnsPropagationTimeout
}

func (d *DnsUserProfile) GetDnsPollingInterval() int64 {
	return d.DnsPollingInterval
}

func (d *DnsUserProfile) GetCertBundle() bool {
	return d.CertBundle
}
