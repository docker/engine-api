// Package transport provides function to send request to remote endpoints.
package transport

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// apiTransport holds information about the http transport to connect with the API.
type apiTransport struct {
	*http.Client
	*tlsInfo
	transport *http.Transport
}

// NewTransportWithHTTP creates a new transport based on the provided proto, address and http client.
// It uses Docker's default http transport configuration if the client is nil.
// It does not modify the client's transport if it's not nil.
func NewTransportWithHTTP(proto, addr string, client *http.Client) (Client, error) {
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

	return &apiTransport{
		Client:    client,
		tlsInfo:   &tlsInfo{transport.TLSClientConfig},
		transport: transport,
	}, nil
}

// CancelRequest stops a request execution.
func (a *apiTransport) CancelRequest(req *http.Request) {
	a.transport.CancelRequest(req)
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
	} else if proto == "npipe" {
		tr.Dial = func(_, _ string) (net.Conn, error) {
			return DialPipe(addr, timeout)
		}
	} else {
		tr.Proxy = http.ProxyFromEnvironment
		tr.Dial = (&net.Dialer{Timeout: timeout}).Dial
	}

	return tr
}

var _ Client = &apiTransport{}
