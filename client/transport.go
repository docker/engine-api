package client

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

// apiTransport holds information about the http transport to connect with the API.
type apiTransport struct {
	// httpClient holds the client transport instance. Exported to keep the old code running.
	httpClient *http.Client
	// scheme holds the scheme of the client i.e. https.
	scheme string
	// tlsConfig holds the tls configuration to use in hijacked requests.
	tlsConfig *tls.Config
}

// newAPITransport creates a new transport based on the provided proto, address and client.
// It uses Docker's default http transport configuration if the client is nil.
// It does not modify the client's transport if it's not nil.
func newAPITransport(proto, addr string, client *http.Client) (*apiTransport, error) {
	scheme := "http"
	var transport *http.Transport

	if client != nil {
		tr, ok := client.Transport.(*http.Transport)
		if !ok {
			return nil, fmt.Errorf("unable to verify TLS configuration, invalid transport %v", client.Transport)
		}
		transport = tr
	} else {
		transport = defaultTransport(proto, addr)
		client = &http.Client{
			Transport: transport,
		}
	}

	if transport.TLSClientConfig != nil {
		scheme = "https"
	}

	return &apiTransport{
		httpClient: client,
		scheme:     scheme,
		tlsConfig:  transport.TLSClientConfig,
	}, nil
}

// HTTPClient returns the http client.
func (a *apiTransport) HTTPClient() *http.Client {
	return a.httpClient
}

// Scheme returns the api scheme.
func (a *apiTransport) Scheme() string {
	return a.scheme
}

// TLSConfig returns the TLS configuration.
func (a *apiTransport) TLSConfig() *tls.Config {
	return a.tlsConfig
}

// IsTLS returns true if there is a TLS configuration.
func (a *apiTransport) IsTLS() bool {
	return a.tlsConfig != nil
}

// defaultTransport creates a new http.Transport with Docker's
// default transport configuration.
func defaultTransport(proto, addr string) *http.Transport {
	tr := new(http.Transport)

	// Why 32? See https://github.com/docker/docker/pull/8035.
	timeout := 32 * time.Second
	if proto == "unix" {
		// No need for compression in local communications.
		tr.DisableCompression = true
		tr.Dial = func(_, _ string) (net.Conn, error) {
			return net.DialTimeout(proto, addr, timeout)
		}
	} else {
		tr.Proxy = http.ProxyFromEnvironment
		tr.Dial = (&net.Dialer{Timeout: timeout}).Dial
	}

	return tr
}
