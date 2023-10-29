package certark

// storage latest state for a task
type StateProfile struct {
	TaskName   string   `json:"task_name"`
	RequestDir string   `json:"request_dir"`
	Cert       string   `json:"cert"`       // latest certificate
	Priv       string   `json:"priv"`       // latest private key
	Domains    []string `json:"domains"`    // latest domains
	Counter    int      `json:"counter"`    // request counter
	LastState  string   `json:"last_state"` // last request state file
	RequestAt  string   `json:"request_at"`
	UrlCheckAt string   `json:"url_check_at"`
}

// storage task running history
type StateHistory struct {
	TaskName   string   `json:"task_name"`
	RequestDir string   `json:"request_dir"`
	Cert       string   `json:"cert"`    // certificate
	Priv       string   `json:"priv"`    // private key
	Domains    []string `json:"domains"` // domains
	Triggger   string   `json:"trigger"` // running reason
	RequestAt  string   `json:"request_at"`
}
