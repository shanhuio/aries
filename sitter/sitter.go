package sitter

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"shanhu.io/misc/rand"
)

type sitter struct {
	config    *Config
	internal1 *oneLine
	internal2 *oneLine
	redirect  *redirect

	old, current *space
	spaces       *spaces

	nextAddr string
	usedAddr string

	spaceToRun chan string
}

func listenAndServe(addr string, h http.Handler) {
	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatal(err)
	}
}

func run(c *Config) {
	s := &sitter{
		config:    c,
		internal1: newOneLine(c.Line1),
		internal2: newOneLine(c.Line2),
		redirect:  newRedirect(c.Internal1),

		spaces:     loadSpaces(c.StateFile),
		spaceToRun: make(chan string, 5),

		nextAddr: c.External1,
		usedAddr: c.External2,
	}

	running := s.spaces.Running

	// remove all other spaces other than the running one.
	for name, sp := range s.spaces.Spaces {
		if name == running {
			continue
		}
		if err := sp.cleanUp(); err != nil {
			log.Println(err)
		}
	}

	log.Printf("serve on %s", c.Service)
	log.Printf("control on %s", c.Control)

	go listenAndServe(c.Control, s)
	go listenAndServe(c.Internal1, s.internal1)
	go listenAndServe(c.Internal2, s.internal2)
	go listenAndServe(c.Service, s.redirect)

	if running != "" {
		log.Printf("restore: %q", running)
		s.schedule(running)
	}
	s.serveSwitching()
}

func (s *sitter) pointTo(addr string) { s.redirect.setHost(addr) }

func (s *sitter) switchTo(name string) (*space, error) {
	switch name {
	case spaceInternal1:
		s.pointTo(s.config.Internal1)
		return nil, nil
	case spaceInternal2:
		s.pointTo(s.config.Internal2)
		return nil, nil
	}
	next, err := s.spaces.get(name)
	if err != nil {
		return nil, err
	}

	nextAddr := s.nextAddr
	if err := next.startAtAddr(nextAddr); err != nil {
		log.Println(err)
	}

	// TODO(h8liu): monitor the health status
	time.Sleep(time.Second * 3)

	s.pointTo(s.nextAddr)
	s.nextAddr, s.usedAddr = s.usedAddr, s.nextAddr
	return next, nil
}

func (s *sitter) save() {
	if err := s.spaces.save(); err != nil {
		log.Println(err)
	}
}

// serveSwitching is the background thread for controlling
// the redirection and space life cycle.
func (s *sitter) serveSwitching() {
	for spaceName := range s.spaceToRun {
		log.Printf("switching to %q", spaceName)
		next, err := s.switchTo(spaceName)
		if err != nil {
			log.Println(err)
			continue
		}

		old := s.current
		s.current = next
		s.spaces.setRunning(next.Name)
		s.save()

		// turn down the old space
		if old != nil {
			if err := old.stop(); err != nil {
				log.Println(err)
			}
			if err := old.cleanUp(); err != nil {
				log.Println(err)
			}
			s.spaces.remove(old.Name)
			s.save()
		}
	}
}

func (s *sitter) schedule(name string) {
	s.spaceToRun <- name
}

func (s *sitter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch path {
	case "/i1":
		s.schedule(spaceInternal1)
	case "/i2":
		s.schedule(spaceInternal2)
	case "/new":
		name := rand.HexBytes(6)
		dir := filepath.Join(s.config.Home, name)
		sp, err := newSpace(name, dir)
		if replyError(w, err) {
			return
		}
		if err := s.spaces.add(sp); replyError(w, err) {
			sp.cleanUp()
			return
		}
		s.save()
		fmt.Fprintln(w, name)
	case "/put":
		query := r.URL.Query()
		space := query.Get("space")
		sp, err := s.spaces.get(space)
		if replyError(w, err) {
			return
		}
		if err := unzip(r.Body, sp.Dir); replyError(w, err) {
			return
		}
	case "/run":
		query := r.URL.Query()
		bin := query.Get("bin")
		space := query.Get("space")
		if err := s.spaces.setBin(space, bin); replyError(w, err) {
			return
		}
		s.save()

		s.schedule(space)
	default:
		log.Println(path)
	}
}

// ServeConfig runs the sitter and starts listening on all the addresses in the
// config structure.
func ServeConfig(c *Config) error {
	run(c)
	panic("unreachable")
}

// Serve runs the default sitter.
func Serve(local bool, stateFile string) error {
	c := Default(local)
	c.StateFile = stateFile
	return ServeConfig(c)
}
