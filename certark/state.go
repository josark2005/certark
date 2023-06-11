package certark

type StateProfile struct {
	TaskName   string   `json:"task_name"`
	RequestDir string   `json:"request_dir"`
	Cert       string   `json:"cert"`
	Priv       string   `json:"priv"`
	Domain     []string `json:"domain"`
	Acme       string   `json:"acmue"`
	RequestAt  string   `json:"request_at"`
	UrlCheckAt string   `json:"url_check_at"`
}
