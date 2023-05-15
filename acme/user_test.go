package acme_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/registration"
	"github.com/jokin1999/certark/acme"
)

func TestGenprivateKey(t *testing.T) {
	pk, err := acme.GenPrivateKey()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(pk)
	}
}

func TestGetCert(t *testing.T) {
	// pk, err := acme.GenPrivateKey()
	// if err != nil {
	// 	panic(err)
	// }

	// log.Println(acme.PrivateKeyEncode(pk))

	domain := "test3.msjnb.eu.org"
	authEmail := "327928971@qq.com"
	authToken := "C7zehIHFwO1QwhKZg-NTy3zW4GaZqStna-nxTErj"
	acmeUser := "test@inspc.ml"
	acmeUserPrivateKey := `
-----BEGIN PRIVATE KEY-----
MHcCAQEEIJfngFl6/666X/3qZ62ME/tvndABOoR0Xrk9tzfjG9e3oAoGCCqGSM49
AwEHoUQDQgAEN2UuJAZODWWuDSKc6qXZCGJ1eQF6NVW+Lh3DJURcFuG58MYM2nVu
JKasivJ7P4A7IdQH2mV2YogyVXxU2ifHcg==
-----END PRIVATE KEY-----`

	u := acme.AcmeUser{
		Email: acmeUser,
		Key:   acme.PrivateKeyDecode(acmeUserPrivateKey),
	}

	config := lego.NewConfig(&u)
	config.CADirURL = lego.LEDirectoryStaging
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		panic(err)
	}

	cfconf := cloudflare.Config{
		AuthEmail:          authEmail,
		AuthToken:          authToken,
		TTL:                120,
		PropagationTimeout: 1 * time.Minute,
		PollingInterval:    5 * time.Second,
	}
	cf, err := cloudflare.NewDNSProviderConfig(&cfconf)
	if err != nil {
		panic(err)
	}

	err = client.Challenge.SetDNS01Provider(cf)
	if err != nil {
		panic(err)
	}

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}
	u.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(certificates.PrivateKey))
	fmt.Println(string(certificates.Certificate))
}
