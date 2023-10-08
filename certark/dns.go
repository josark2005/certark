package certark

type DNSProfile struct {
	Enabled bool `json:"enabled"`

	DNSProvider  string `json:"dns_provider"`
	DNSAuthUser  string `json:"dns_authuser"`
	DNSAuthKey   string `json:"dns_authkey"`
	DNSAuthToken string `json:"dns_authtoken"`
	DNSZoneToken string `json:"dns_zonetoken"`
}
