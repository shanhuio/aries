package main

import (
	"flag"
	"io/ioutil"
	"log"

	"shanhu.io/aries/creds"
)

func ne(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	out := flag.String("out", "", "key path to output")
	nopass := flag.Bool("nopass", false, "no passphrase")
	flag.Parse()

	var pwd []byte
	var err error

	if !*nopass {
		pwd, err = creds.ReadPassword("Key passphrase: ")
		ne(err)
	}

	pri, pub, err := creds.GenerateKey(pwd)
	ne(err)

	if *out == "" {
		*out, err = creds.HomeFile("key")
		ne(err)
	}

	ne(ioutil.WriteFile(*out+".pem", pri, 0600))
	ne(ioutil.WriteFile(*out+".pub", pub, 0600))
}
