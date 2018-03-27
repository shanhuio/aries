package creds

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func pemBlock(k *rsa.PrivateKey, pwd []byte) (*pem.Block, error) {
	const pemType = "RSA PRIVATE KEY"

	if pwd == nil {
		return &pem.Block{
			Type:  pemType,
			Bytes: x509.MarshalPKCS1PrivateKey(k),
		}, nil
	}

	return x509.EncryptPEMBlock(
		rand.Reader, pemType,
		x509.MarshalPKCS1PrivateKey(k),
		pwd, x509.PEMCipherDES,
	)
}

// GenerateKey generates a private/public key pair with the given passphrase.
func GenerateKey(passphrase []byte, n int) (pri, pub []byte, err error) {
	key, err := rsa.GenerateKey(rand.Reader, n)
	if err != nil {
		return nil, nil, err
	}

	b, err := pemBlock(key, passphrase)
	if err != nil {
		return nil, nil, err
	}

	priBuf := new(bytes.Buffer)
	if err := pem.Encode(priBuf, b); err != nil {
		return nil, nil, err
	}
	pubKey, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	pub = ssh.MarshalAuthorizedKey(pubKey)
	return priBuf.Bytes(), pub, nil
}

// ReadPrivateKey reads a private key from a key file.
func ReadPrivateKey(pemFile string, tty bool) (*rsa.PrivateKey, error) {
	bs, err := ReadPrivateFile(pemFile)
	if err != nil {
		return nil, err
	}
	b, _ := pem.Decode(bs)
	if b == nil {
		return nil, fmt.Errorf("%q decode failed", pemFile)
	}

	if !x509.IsEncryptedPEMBlock(b) {
		return x509.ParsePKCS1PrivateKey(b.Bytes)
	}

	if !tty {
		return nil, fmt.Errorf("%q is encrypted", pemFile)
	}

	prompt := fmt.Sprintf("Passphrase for %s: ", pemFile)
	pwd, err := ReadPassword(prompt)
	if err != nil {
		return nil, err
	}

	der, err := x509.DecryptPEMBlock(b, pwd)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PrivateKey(der)
}
