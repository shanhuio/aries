package ariestest

import (
	"sort"

	"shanhu.io/aries"
	"shanhu.io/aries/https/httpstest"
)

// HTTPSServers creates an HTTPS server that serves a set of
// domains.
func HTTPSServers(sites map[string]aries.Service) (*httpstest.Server, error) {
	m := aries.NewHostMux()
	var domains []string
	for domain, s := range sites {
		m.Set(domain, s)
		domains = append(domains, domain)
	}
	sort.Strings(domains)

	return httpstest.NewServer(domains, aries.Serve(m))
}

// HTTPSServer creates an HTTPS server that serves for a domain.
func HTTPSServer(domain string, s aries.Service) (*httpstest.Server, error) {
	return HTTPSServers(map[string]aries.Service{
		domain: s,
	})
}
