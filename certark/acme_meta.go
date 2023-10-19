package certark

type AcmeUserProfile struct {
	Email      string `json:"email"`
	PrivateKey string `json:"privatekey"`
	Enabled    bool   `json:"enabled"`
}

var DefaultAcmeUserProfile = AcmeUserProfile{
	Email:      "",
	PrivateKey: "",
	Enabled:    true,
}
