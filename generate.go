package main

import (
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/urfave/cli/v2"
)

func humanTime(now, when time.Time) string {
	return fmt.Sprintf("%s (%s)", when, humanize.RelTime(when, now, "ago", "from now"))
}

func Generate(c *cli.Context) error {
	if c.Int64("not-before") == 0 || c.Int64("not-after") == 0 {
		fmt.Printf(`The --not-before and --not-after flag need to be set
to the start and end time of the CA's validity.

The time must be provided in POSIX / UNIX time, the number of seconds since
January 1st, 1970. If this is a thing you are uncomfortable doing, no worries!

Just head to a website like https://www.epochconverter.com to work out the right
numbers to put in there.

Example usage:

  $ mockca generate --bits 2048 --not-before 1512483689 --not-after 1544019689

`)
		return nil
	}

	notBefore := time.Unix(c.Int64("not-before"), 0)
	notAfter := time.Unix(c.Int64("not-after"), 0)
	now := time.Now()

	path := c.String("root")
	bits := c.Int("bits")

	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}

	ca, err := GenerateCA(path, bits, notBefore, notAfter)
	if err != nil {
		return err
	}

	if err := ca.Save(); err != nil {
		return err
	}

	fmt.Printf("CA Serial:             %x\n", ca.Cert.SerialNumber)
	fmt.Printf("CA Subject CommonName: %s\n", ca.Cert.Subject.CommonName)
	fmt.Printf("Starting Time:         %s\n", humanTime(now, ca.Cert.NotBefore))
	fmt.Printf("Ending Time:           %s\n", humanTime(now, ca.Cert.NotAfter))
	fmt.Printf("\n")

	if err := pem.Encode(os.Stdout, &pem.Block{Type: "CERTIFICATE", Bytes: ca.Cert.Raw}); err != nil {
		return err
	}

	return nil
}
