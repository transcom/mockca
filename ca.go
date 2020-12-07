package main

import (
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
)

// High level encapsulation of some of the CA basics. This contains the root
// directory, the loaded Private Key, and the loaded x509 Certificate.
type CA struct {
	Root    string
	PrivKey *rsa.PrivateKey
	Cert    *x509.Certificate
}

// Create a new DoD-Like x509 Certificate from a person, and a CSR that they
// send us.
//
// Most importantly, this will disregard the CSR's data entirely in favor of
// generating the data from what we're supposed to know about.
func (ca CA) Issue(who Employee, csr *x509.CertificateRequest) ([]byte, error) {
	template, err := who.CertificateTemplate(ca.Cert.NotBefore, ca.Cert.NotAfter)
	if err != nil {
		return nil, err
	}
	template.PublicKey = csr.PublicKey
	template.PublicKeyAlgorithm = csr.PublicKeyAlgorithm
	template.SignatureAlgorithm = x509.SHA256WithRSA
	return x509.CreateCertificate(rand.Reader, template, ca.Cert, csr.PublicKey, ca.PrivKey)
}

// load the rsa private key from the filesystem (eww, gross)
func (ca *CA) loadKey() error {
	data, _, err := ca.load("key.pem")
	if err != nil {
		return err
	}
	ca.PrivKey, err = x509.ParsePKCS1PrivateKey(data)
	return err
}

// Load the x509 Certificate from the filesystem
func (ca *CA) loadCert() error {
	data, _, err := ca.load("ca.crt")
	if err != nil {
		return err
	}
	ca.Cert, err = x509.ParseCertificate(data)
	return err
}

// Internal function to pull in the rsa private key and the x509 Certificate
// for use later on. This needs the pointer reciever to mutate the CA
// object.
func (ca *CA) init() error {
	if err := ca.loadKey(); err != nil {
		return err
	}
	return ca.loadCert()
}

// Load the bits and bobs off the filesystem.
func (ca CA) load(name string) ([]byte, string, error) {
	path := filepath.Join(ca.Root, name)
	fd, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, "", err
	}
	block, _ := pem.Decode(data)
	if err != nil {
		return nil, "", err
	}
	return block.Bytes, block.Type, nil
}

// Save some bytes to the filesystem.
//
// This will write a PEM blob at root/name, with the Type of V, and
// contents in Data.
func (ca CA) save(name, v string, data []byte) error {
	path := filepath.Join(ca.Root, name)

	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()

	return pem.Encode(fd, &pem.Block{Type: v, Bytes: data})
}

// Write out the rsa Private Key.
func (ca CA) saveKey() error {
	return ca.save("key.pem", "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(ca.PrivKey))
}

// Write out the x509 Certificate.
func (ca CA) saveCert() error {
	return ca.save("ca.crt", "CERTIFICATE", ca.Cert.Raw)
}

// Save all the internal bits and bobs to the filesystem.
func (ca CA) Save() error {
	if err := ca.saveKey(); err != nil {
		return err
	}
	return ca.saveCert()
}

// Load the private key and Certificate off the Filesystem at the following
// location.
func LoadCA(path string) (*CA, error) {
	ca := CA{Root: path}
	if err := ca.init(); err != nil {
		return nil, err
	}
	return &ca, nil
}

// Generate a new RSA Private Key, then write out the bits we generated out to
// the Filesystem.
func GenerateCA(path string, bits int, notBefore, notAfter time.Time) (*CA, error) {
	template, err := certificateTemplate(pkix.Name{
		Country:            []string{"US"},
		Organization:       []string{"U.S. Government"},
		OrganizationalUnit: []string{"DoD", "PKI"},
		CommonName:         "DOD EMAIL CA-NaN",
	}, notBefore, notAfter)
	if err != nil {
		return nil, err
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign
	template.KeyUsage |= x509.KeyUsageKeyEncipherment
	template.KeyUsage |= x509.KeyUsageKeyAgreement

	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, priv.Public(), priv)
	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, err
	}

	return &CA{
		Root:    path,
		Cert:    cert,
		PrivKey: priv,
	}, nil

}
