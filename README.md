mockca
======

Quick Start
-----

1. Run `go build -o mockca-build`
2. Generate CA files given start and end timestamps: `./mockca-build generate -not-before 1512436787 -not-after 1513436787`
3. Add identities to `identities.csv` in form `First,M,Last,EDIPI`
4. Run `python generate.py`. 
5. Certificates will be generated in the `certs` folder. The password is set by default to `1234`.

Usage
-----

First, generate the Certificate Authority files using the `-not-before` (beginning) and `-not-after` (expiry) flags.

Consult an [epoch time converter](https://www.epochconverter.com) to get the POSIX / UNIX time values required.

```bash
./mockca generate -not-before 1512436787 -not-after 1513436787
```

Somewhere else, on another computer, a user should create their CSR
and send the file `me.csr` to the person running the CA.
-----
`openssl req -new -nodes -out me.csr -subj '/CN=new cert/'`

Later, we can take a CSR, and create a new Certificate that looks right only using the CSR Public Key.
-----

```bash
./mockca sign \
    -first-name Dustin \
    -last-name Ellingson \
    -middle-name B \
    -dod-id 1277111111 \
    -email dustin.elligson.1@us.af.mil \
    -org usaf \
    me.csr
```
