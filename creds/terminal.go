package creds

import (
	"fmt"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// ReadPassword reads a password from stdin.
func ReadPassword(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	line, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Println()
	if err != nil {
		return nil, err
	}
	return []byte(line), nil
}
