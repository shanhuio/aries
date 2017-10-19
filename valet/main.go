package valet

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"shanhu.io/aries"
)

func ne(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func makeProxy(c *Config) *httputil.ReverseProxy {
	proxy := new(httputil.ReverseProxy)
	proxy.FlushInterval = time.Second * 3
	proxy.Director = func(req *http.Request) {
		host := req.Host

		url := req.URL
		url.Scheme = "http"

		if host == "" {
			log.Println("empty host")
			return
		}

		if host == c.Control {
			url.Host = controlHost
			return
		}

		if mapTo, ok := c.Hosts[host]; ok {
			url.Host = mapTo
			return
		}

		// because we are white listing, this should not happen
		log.Printf("unexpected host %q", host)
		url.Host = "localhost:8000"
	}
	return proxy
}

func redirectHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS != nil || r.Host == "" {
		http.Error(w, "not found", 404)
	}

	u := r.URL
	u.Host = r.Host
	u.Scheme = "https"
	http.Redirect(w, r, u.String(), 302)
}

const controlHost = "localhost:8001"

func runHTTPServer(addr string, h http.Handler) {
	s := &http.Server{
		Addr:    addr,
		Handler: h,
	}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()
}

func runHTTPServerFunc(addr string, f http.HandlerFunc) {
	runHTTPServer(addr, f)
}

func serve(c *Config) {
	policy := func(_ context.Context, host string) error {
		if !c.hasHost(host) {
			return fmt.Errorf("%q is not in the whitelist", host)
		}
		return nil
	}

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: policy,
		Cache:      autocert.DirCache(c.cache()),
	}

	runHTTPServerFunc(":http", redirectHTTP)
	runHTTPServer(controlHost, &aries.Handler{
		Func:  makeControl(c),
		HTTPS: true,
	})

	proxy := makeProxy(c)
	s := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
		},
		Handler: proxy,
	}
	log.Fatal(s.ListenAndServeTLS("", ""))
}

// Main is the main entrance of the smlfront binary.
func Main() {
	var (
		config = flag.String("config", "config.json", "JSON config file")
	)
	c, err := loadConfig(*config)
	ne(err)

	serve(c)
}
