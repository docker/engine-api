package transport

import (
	"crypto/tls"
	"net"
	"net/http"
)

// Dialer is an interface thet clients must implement
// to be able to set up hijacked connections.
type Dialer interface {
	Dial(string, string) (net.Conn, error)
}

// Sender is an interface that clients must implement
// to be able to send requests to a remote connection.
type Sender interface {
	// Do sends request to a remote endpoint.
	Do(*http.Request) (*http.Response, error)
}

// Client is an interface that abstracts all remote connections.
type Client interface {
	Dialer
	Sender
	// Secure tells whether the connection is secure or not.
	Secure() bool
	// Scheme returns the connection protocol the client uses.
	Scheme() string
	// TLSConfig returns any TLS configuration the client uses.
	TLSConfig() *tls.Config
}

// tlsInfo returns information about the TLS configuration.
type tlsInfo struct {
	tlsConfig *tls.Config
}

// TLSConfig returns the TLS configuration.
func (t *tlsInfo) TLSConfig() *tls.Config {
	return t.tlsConfig
}

// Scheme returns protocol scheme to use.
func (t *tlsInfo) Scheme() string {
	if t.tlsConfig != nil {
		return "https"
	}
	return "http"
}

// Secure returns true if there is a TLS configuration.
func (t *tlsInfo) Secure() bool {
	return t.tlsConfig != nil
}
