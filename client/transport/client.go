package transport

import (
	"crypto/tls"
	"net/http"
)

// Client is an interface that abstracts all remote connections.
type Client interface {
	// Do sends request to a remote endpoint.
	Do(req *http.Request) (resp *http.Response, err error)
	// Secure tells whether the connection is secure or not.
	Secure() bool
	// Scheme returns the connection protocol the client uses.
	Scheme() string
	// TLSConfig returns any TLS configuration the client uses.
	TLSConfig() *tls.Config
}
