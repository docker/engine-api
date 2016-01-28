package transport

import (
	"crypto/tls"
	"net/http"
)

// Client is an interface that abstracts all remote connections.
type Client interface {
	// Do sends request to a remote endpoint.
	Do(*http.Request) (*http.Response, error)
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
