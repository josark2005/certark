package acme

type AcmeProfile interface {
	GetEmail() string
	GetPrivateKey() string
}
