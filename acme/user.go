package acme

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/jokin1999/certark/ark"
)

const (
	MODE_STAGING    = 0
	MODE_PRODUCTION = 1
)

type AcmeUser struct {
	Email        string
	Registration *registration.Resource
	Key          crypto.PrivateKey
}

func (u *AcmeUser) GetEmail() string {
	return u.Email
}
func (u AcmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}

func GenPrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// encode ecdsa private key
func PrivateKeyEncode(privateKey *ecdsa.PrivateKey) string {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	return string(pemEncoded)
}

// decode pem to ecdsa private key
func PrivateKeyDecode(pemEncoded string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	return privateKey
}

func RegisterAcmeUser(email string, mode int) string {
	// generate user info
	privateKey, _ := GenPrivateKey()
	u := AcmeUser{
		Email: email,
		Key:   privateKey,
	}

	// generate lego config
	config := lego.NewConfig(&u)

	// set server
	if mode == MODE_STAGING {
		config.CADirURL = lego.LEDirectoryStaging
	} else {
		config.CADirURL = lego.LEDirectoryProduction
	}
	config.Certificate.KeyType = certcrypto.RSA2048
	client, err := lego.NewClient(config)
	if err != nil {
		ark.Error().Err(err).Msg("Register acme user failed")
	}

	// register user
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		ark.Error().Err(err).Msg("Register acme user failed")
	}
	u.Registration = reg

	return PrivateKeyEncode(privateKey)
}
