package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"shanhu.io/aries/creds"
	"shanhu.io/misc/osutil"
)

func ne(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	out := flag.String("out", "", "key path to output")
	nopass := flag.Bool("nopass", false, "no passphrase")
	nbit := flag.Int("nbit", 4096, "number of RSA bits")
	flag.Parse()

	var pwd []byte
	var err error

	if !*nopass {
		pwd, err = creds.ReadPassword("Key passphrase: ")
		ne(err)
	}

	pri, pub, err := creds.GenerateKey(pwd, *nbit)
	ne(err)

	if *out == "" {
		*out, err = creds.HomeFile("key")
		ne(err)
	}

	pemPath := *out + ".pem"
	if yes, err := osutil.Exist(pemPath); err != nil {
		ne(err)
	} else if yes {
		ne(fmt.Errorf("key file %q already exists", pemPath))
	}

	ne(ioutil.WriteFile(*out+".pem", pri, 0600))
	ne(ioutil.WriteFile(*out+".pub", pub, 0600))
}
