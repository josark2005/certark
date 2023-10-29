package certificate_test

import (
	"testing"

	"github.com/jokin1999/certark/certark/certificate"
)

const TestDomain = "josark.com"
const TestPort = 443

func TestFetchOnlineCertificate(t *testing.T) {
	certs, err := certificate.FetchOnlineCertificate(TestDomain, TestPort)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(certs[0].DNSNames)
		t.Log(certs[0].BasicConstraintsValid)
		t.Log(certs[0].NotBefore)
		t.Log(certs[0].NotAfter)
	}
}

func TestCheckCertificatValidSimple(t *testing.T) {
	certs, err := certificate.FetchOnlineCertificate(TestDomain, TestPort)
	if err != nil {
		t.Error(err)
	}
	t.Log(certificate.CheckCertificatValidSimple("josark.com", certs[0]))
}
