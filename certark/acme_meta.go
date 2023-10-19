package certark

type AcmeUserProfile struct {
	Email      string `json:"email"`
	PrivateKey string `json:"privatekey"`
	Enabled    bool   `json:"enabled"`
}
