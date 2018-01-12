package httpstest

import (
	"context"
	"net"
)

func sink(sinkAddr string) func(
	ctx context.Context, net, addr string,
) (net.Conn, error) {
	d := new(net.Dialer)
	return func(ctx context.Context, net, addr string) (net.Conn, error) {
		return d.DialContext(ctx, net, sinkAddr)
	}
}

func sinkHTTPS(httpAddr, httpsAddr string) func(
	ctx context.Context, net, addr string,
) (net.Conn, error) {
	d := new(net.Dialer)
	return func(ctx context.Context, netStr, addr string) (net.Conn, error) {
		_, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		sinkAddr := httpAddr
		if port == "443" || port == "https" {
			sinkAddr = httpsAddr
		}
		return d.DialContext(ctx, netStr, sinkAddr)
	}
}
