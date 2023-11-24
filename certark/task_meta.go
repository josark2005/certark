package certark

type TaskProfile struct {
	TaskName string   `json:"task_name"`
	Domains   []string `json:"domain"`
	AcmeUser string   `json:"acme_user"`
	Enabled  bool     `json:"enabled"`

	DnsProfile string `json:"dns_profile"`

	UrlCheckEnable   bool   `json:"url_check_enable"`
	UrlCheckTarget   string `json:"url_check_target"`
	UrlCheckInterval int64  `json:"url_check_interval"` // in day, 1 is recommanded

	CertExportPath        string `json:"cert_export_path"`         // export cert after cert updating
	PostCertUpdateCommand string `json:"post_cert_update_command"` // command runs after cert udpating
}

var DefaultTaskProfile = TaskProfile{
	TaskName: "default",
	Domains:   []string{},
	AcmeUser: "",
	Enabled:  true,

	DnsProfile: "",

	UrlCheckEnable:   false,
	UrlCheckTarget:   "",
	UrlCheckInterval: 1,

	CertExportPath:        "",
	PostCertUpdateCommand: "",
}
