package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

const path string = "./src/certs/"

func makeCA(subject *pkix.Name) (*x509.Certificate, *rsa.PrivateKey, error) {
	// creating a CA which will be used to sign all of our certificates using the x509 package from the Go Standard Library
	caCert := &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               *subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10*365, 0, 0),
		IsCA:                  true, // <- indicating this certificate is a CA certificate.
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	// generate a private key for the CA
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("Generate the CA Private Key error: %v\n", err)
		return nil, nil, err
	}

	// create the CA certificate
	caBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, &caKey.PublicKey, caKey)
	if err != nil {
		log.Printf("Create the CA Certificate error: %v\n", err)
		return nil, nil, err
	}

	// Create the CA PEM files
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	if err := os.WriteFile(path+"ca.crt", caPEM.Bytes(), 0644); err != nil {
		log.Printf("Write the CA certificate file error: %v\n", err)
		return nil, nil, err
	}

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caKey),
	})
	if err := os.WriteFile(path+"ca.key", caPrivKeyPEM.Bytes(), 0644); err != nil {
		log.Printf("Write the CA private key file error: %v\n", err)
		return nil, nil, err
	}
	return caCert, caKey, nil
}

func makeCert(caCert *x509.Certificate, caKey *rsa.PrivateKey, subject *pkix.Name, name string, DNSNamesArray []string) error {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject:      *subject,
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
		DNSNames:     DNSNamesArray,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Printf("Generate the Key error: %v\n", err)
		return err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, &certKey.PublicKey, caKey)
	if err != nil {
		log.Printf("Generate the certificate error: %v\n", err)
		return err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err := os.WriteFile(path+name+".crt", certPEM.Bytes(), 0644); err != nil {
		log.Printf("Write the CA certificate file error: %v\n", err)
		return err
	}

	certKeyPEM := new(bytes.Buffer)
	pem.Encode(certKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certKey),
	})
	if err := os.WriteFile(path+name+".key", certKeyPEM.Bytes(), 0644); err != nil {
		log.Printf("Write the CA certificate file error: %v\n", err)
		return err
	}
	return nil
}

func main() {
	// Check if at least one domain is provided as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_cert.go <domain1> [domain2] [domain3] ...")
		return
	}

	domains := os.Args[1:]

	// Create the DNSNamesArray including "localhost" and the provided domains
	DNSNamesArray := append([]string{"localhost"}, domains...)

	fmt.Println("Generating certificates for domains:", DNSNamesArray)
	subject := pkix.Name{
		Country:            []string{"IN"},
		Organization:       []string{"DBOM"},
		OrganizationalUnit: []string{"DBOM"},
		Locality:           []string{"BLR"},
		Province:           []string{"KA"},
		CommonName:         "*",
		Names:              []pkix.AttributeTypeAndValue{},
		ExtraNames:         []pkix.AttributeTypeAndValue{},
	}
	caCert, caKey, err := makeCA(&subject)
	if err != nil {
		log.Fatalf("make CA Certificate error!")
	}
	log.Println("Created the CA certificate successfully.")

	for _, domain := range os.Args[1:] {
		if err := makeCert(caCert, caKey, &subject, domain, DNSNamesArray); err != nil {
			log.Printf("Failed to make certificate for domain %s: %v", domain, err)
		} else {
			log.Printf("Created and signed the certificate for %s successfully.", domain)
		}
	}
}
