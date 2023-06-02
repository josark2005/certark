package certark

type TaskProfile struct {
	TaskName string   `json:"task_name"`
	Domain   []string `json:"domain"`
	AcmeUser string   `json:"acme_user"`
	Enabled  bool     `json:"enabled"`

	DNSProvider           string `json:"dns_provider"`
	DNSAuthUser           string `json:"dns_authuser"`
	DNSAuthKey            string `json:"dns_authkey"`
	DNSAuthToken          string `json:"dns_authtoken"`
	DNSZoneToken          string `json:"dns_zonetoken"`
	DNSTTL                int64  `json:"dns_ttl"`                 // ttl 120 is recommanded
	DNSPropagationTimeout int64  `json:"dns_propagation_timeout"` // in millisecond, 60*1000 is recommanded
	DNSPollingInterval    int64  `json:"dns_polling_interval"`    // in millisecond, 5 *1000 is recommanded

	UrlCheckEnable   bool   `json:"url_check_enable"`
	UrlCheckTarget   string `json:"url_check_target"`
	UrlCheckInterval int64  `json:"url_check_interval"` // in day, 1 is recommanded
}

var DefaultTaskProfile = TaskProfile{
	TaskName:              "default",
	Domain:                []string{},
	AcmeUser:              "",
	Enabled:               true,
	DNSProvider:           "",
	DNSAuthUser:           "",
	DNSAuthKey:            "",
	DNSAuthToken:          "",
	DNSZoneToken:          "",
	DNSTTL:                120,
	DNSPropagationTimeout: 60,
	DNSPollingInterval:    5,
	UrlCheckEnable:        false,
	UrlCheckTarget:        "",
	UrlCheckInterval:      1,
}
