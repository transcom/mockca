package main

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"crypto/x509"

	"github.com/urfave/cli"
)

func LoadEmployee(c *cli.Context) (*Employee, error) {
	var org Organization
	switch c.String("org") {
	case "other":
		org = Other
	case "osd":
		org = SecDef
	case "usaf":
		org = AirForce
	case "army":
		org = Army
	case "navy":
		org = Navy
	case "marines":
		org = Marines
	default:
		return nil, fmt.Errorf("mockca: unknown org")
	}

	return &Employee{
		FirstName:    c.String("first-name"),
		MiddleName:   c.String("middle-name"),
		LastName:     c.String("last-name"),
		Email:        c.String("email"),
		DODID:        c.String("dod-id"),
		Organization: org,
	}, nil
}

func Sign(c *cli.Context) error {
	path := c.GlobalString("root")
	ca, err := LoadCA(path)
	if err != nil {
		return err
	}

	employee, err := LoadEmployee(c)
	if err != nil {
		return err
	}

	for _, el := range c.Args() {
		fd, err := os.Open(el)
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(fd)
		if err != nil {
			return err
		}
		block, _ := pem.Decode(data)
		csr, err := x509.ParseCertificateRequest(block.Bytes)
		if err != nil {
			return err
		}

		der, err := ca.Issue(*employee, csr)
		if err != nil {
			return err
		}

		pem.Encode(os.Stdout, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: der,
		})
	}

	return nil
}
