package static

import (
	"flag"
	"log"
	"net/http"

	"shanhu.io/misc/jsonfile"
)

// Config contains the config file for the smlstatic binary.
type Config struct {
	Dir string // Home directory
}

func ne(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Main is the main entrance for smlstatic binary
func Main() {
	addr := flag.String("addr", "localhost:8000", "listen address")
	config := flag.String("config", "config.json", "config file path")

	flag.Parse()

	var c Config
	ne(jsonfile.Read(*config, &c))

	log.Printf("listening at %s", *addr)
	http.Handle("/", http.FileServer(http.Dir(c.Dir)))
	ne(http.ListenAndServe(*addr, nil))
}
