package certark

type StateProfile struct {
	Task       string   `json:"task"`
	Request    string   `json:"request"`
	Cert       string   `json:"cert"`
	Priv       string   `json:"priv"`
	Domain     []string `json:"domain"`
	Acme       string   `json:"acmue"`
	RequestAt  string   `json:"request_at"`
	UrlCheckAt string   `json:"url_check_at"`
}
