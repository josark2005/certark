package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"strconv"
	"time"
)

func FetchOnlineCertificate(addr string, port int) ([]*x509.Certificate, error) {
	conn, err := tls.Dial("tcp", addr+":"+strconv.Itoa(port), nil)
	if err != nil {
		return []*x509.Certificate{}, err
	}

	certs := conn.ConnectionState().PeerCertificates
	return certs, nil
}

// check certificate basic constraints and time
func CheckCertificatValidSimple(domain string, cert *x509.Certificate) bool {
	if !cert.BasicConstraintsValid {
		return false
	}

	if time.Now().Before(cert.NotBefore) {
		return false
	}

	if time.Now().After(cert.NotAfter) {
		return false
	}

	return true
}
