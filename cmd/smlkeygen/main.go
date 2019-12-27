package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"shanhu.io/aries/creds"
	"shanhu.io/misc/osutil"
)

type config struct {
	nbit         int
	noPassphrase bool
}

func keygen(output string, config *config) error {
	var passphrase []byte
	if !config.noPassphrase {
		pass, err := creds.ReadPassword("Key passphrase: ")
		if err != nil {
			return err
		}
		passphrase = pass
	}

	pri, pub, err := creds.GenerateKey(passphrase, config.nbit)
	if err != nil {
		return err
	}

	if output == "" {
		out, err := creds.HomeFile("key")
		if err != nil {
			return err
		}
		output = out
	}

	pemPath := output + ".pem"
	if yes, err := osutil.Exist(pemPath); err != nil {
		return err
	} else if yes {
		return fmt.Errorf("key file %q already exists", pemPath)
	}

	if err := ioutil.WriteFile(pemPath, pri, 0600); err != nil {
		return err
	}

	return ioutil.WriteFile(output+".pub", pub, 0600)
}

func main() {
	out := flag.String("out", "", "key path to output")
	nopass := flag.Bool("nopass", false, "no passphrase")
	nbit := flag.Int("nbit", 4096, "number of RSA bits")
	flag.Parse()

	conf := &config{
		nbit:         *nbit,
		noPassphrase: *nopass,
	}

	if err := keygen(*out, conf); err != nil {
		log.Fatal(err)
	}
}
