package security

import (
	"log"
	"time"
	"math/big"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
)

type TLSData struct {
}

func CreateCertificate(Country, Organization, OrganizationalUnit string) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Country:            []string{Country},
			Organization:       []string{Organization},
			OrganizationalUnit: []string{OrganizationalUnit},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 1024)
	pub := &priv.PublicKey
	cab, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		log.Println("create ca failed", err)
		return
	}
	caf := "ca.pem"
	log.Println("write to", caf)
	ioutil.WriteFile(caf, cab, 0777)

}
