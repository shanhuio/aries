package sitter

import (
	"flag"
	"log"
)

// Main is the main entrance function of the baby sitter.
func Main() {
	local := flag.Bool("local", false, "listens on local ports")
	state := flag.String("state", "state.json", "file to save the state")
	flag.Parse()
	if err := Serve(*local, *state); err != nil {
		log.Fatal(err)
	}
}
