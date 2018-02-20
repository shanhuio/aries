package aries

import (
	"log"
	"os"
)

// Log is the logger for internal errors.
var Log = log.New(os.Stderr, "", log.LstdFlags)
