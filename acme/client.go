package acme

import (
	"fmt"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/lego"
)

func setConfigDefault(c *lego.Config, dirType string) {
	if dirType == MODE_PROD {
		c.CADirURL = lego.LEDirectoryProduction
	} else {
		c.CADirURL = lego.LEDirectoryStaging
	}
	c.Certificate.KeyType = certcrypto.RSA2048
}

func NewConfig(email, privateKey string, dirType string) (*AcmeUser, *lego.Config) {
	acmeUser := &AcmeUser{
		Email: email,
		Key:   PrivateKeyDecode(privateKey),
	}
	c := lego.NewConfig(acmeUser)
	setConfigDefault(c, dirType)
	return acmeUser, c
}

func NewConfigWithUser(au *AcmeUser, dirType string) *lego.Config {
	c := lego.NewConfig(au)
	setConfigDefault(c, dirType)
	return c
}

func NewConfigWithProfile(acme AcmeProfile, dirType string) (*AcmeUser, *lego.Config) {
	acmeUser := &AcmeUser{
		Email: acme.GetEmail(),
		Key:   PrivateKeyDecode(acme.GetPrivateKey()),
	}
	c := lego.NewConfig(acmeUser)
	setConfigDefault(c, dirType)
	return acmeUser, c
}

func NewClient(c *lego.Config, acmeUser *AcmeUser, provider challenge.Provider) (*lego.Client, error) {
	client, err := lego.NewClient(c)
	if err != nil {
		return client, err
	}

	// select provider
	err = client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		return client, err
	}

	// query register
	reg, err := client.Registration.ResolveAccountByKey()
	if err != nil {
		return client, err
	}
	acmeUser.Registration = reg
	fmt.Println("b", reg.Body)
	fmt.Println("s", reg.Body.Status)

	return client, nil
}
