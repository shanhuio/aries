package aries

// HostMux routes request to different services based on the incoming host.
type HostMux struct {
	m map[string]Service
}

// NewHostMux creates a new host mux.
func NewHostMux() *HostMux {
	return &HostMux{m: make(map[string]Service)}
}

// Set binds a host domain to a particular service.
func (m *HostMux) Set(host string, s Service) {
	m.m[host] = s
}

// Serve serves an incoming request.
func (m *HostMux) Serve(c *C) error {
	host := c.Req.Host
	s, found := m.m[host]
	if !found {
		return Miss
	}
	return s.Serve(c)
}
