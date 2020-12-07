package main

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"strings"
	"time"
)

// OU that an indivudual may belong to.
type Organization string

var (
	// Federal Affililate, etc.
	Other Organization = "OTHER"

	// SecDef / OSD
	SecDef Organization = "OSD"

	// Military defaults
	AirForce Organization = "USAF"
	Army     Organization = "USA"
	Navy     Organization = "USN"
	Marines  Organization = "USMC"
)

// Some random Bureaucrat
type Employee struct {
	FirstName string
	// If MiddleName is set to "", it will be omitted
	MiddleName string
	LastName   string
	Email      string
	// Something like "1234124124"
	DODID        string
	Organization Organization
}

// Generate a pkix.Name for use in the x509 Subject field.
func (e Employee) Subject() pkix.Name {
	return pkix.Name{
		Country:            []string{"US"},
		Organization:       []string{"U.S. Government"},
		OrganizationalUnit: []string{"DoD", "PKI", string(e.Organization)},
		CommonName:         e.CN(),
	}
}

// Generate a CN string that is used in the x509 Subject.
func (e Employee) CN() string {
	// BUREAUCRAT.JOHN.Q.123123123
	names := []string{
		strings.ToUpper(e.LastName),
		strings.ToUpper(e.FirstName),
	}
	if len(e.MiddleName) != 0 {
		names = append(names, strings.ToUpper(e.MiddleName))
	}
	names = append(names, e.DODID)
	return strings.Join(names, ".")
}

// Generate a template x509 Certificate that a user might have.
func (who Employee) CertificateTemplate(notBefore, notAfter time.Time) (*x509.Certificate, error) {
	cert, err := certificateTemplate(who.Subject(), notBefore, notAfter)
	if err != nil {
		return nil, err
	}
	if len(who.Email) != 0 {
		cert.EmailAddresses = []string{who.Email}
	}
	cert.ExtKeyUsage = []x509.ExtKeyUsage{
		x509.ExtKeyUsageClientAuth,
		x509.ExtKeyUsageEmailProtection,
	}

	cert.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageContentCommitment
	cert.IsCA = false

	return cert, nil
}

// Basic outline of a Certificate.
func certificateTemplate(subject pkix.Name, notBefore, notAfter time.Time) (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	// XXX: support UPN dod-id@mil
	return &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
	}, nil
}
