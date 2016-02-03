package authn

import (
	"net/http"

	"github.com/docker/engine-api/client/logger"
)

// BasicAuthCallback is an interface which a caller may provide for obtaining a user
// name and password to use when attempting Basic authentication with a server.
type BasicAuthCallback func(realm string) (user, password string, err error)

// Basic is an AuthResponder that handles basic authentication.
type Basic struct {
	logger             logger.Logger
	callback           BasicAuthCallback
	username, password string
}

// NewBasicAuth creates a Basic auth responder with a callback
// to resolve basic auth credentials.
func NewBasicAuth(callback BasicAuthCallback) AuthResponder {
	return &Basic{
		logger:   logger.Silent{},
		callback: callback,
	}
}

// SetLogger sets the logger for the Basic auth responder.
func (b *Basic) SetLogger(l logger.Logger) {
	b.logger = l
}

// Scheme returns the scheme the Basic auth responder handles.
func (b *Basic) Scheme() string {
	return "Basic"
}

// AuthRespond handles authentication for the Basic auth responder.
func (b *Basic) AuthRespond(challenge string, req *http.Request) (result bool, err error) {
	if b.username != "" && b.password != "" {
		b.logger.Debug("using previously-supplied Basic username and password")
		req.SetBasicAuth(b.username, b.password)
		return true, nil
	}

	if b.callback == nil {
		b.logger.Debug("failed to obtain user name and password for Basic auth")
		return false, nil
	}

	realm, _ := getParameter(challenge, "realm")
	username, password, err := b.callback(realm)
	if err != nil {
		return false, err
	}
	if username == "" {
		b.logger.Debug("failed to obtain user name for Basic auth")
		return false, nil
	}
	if password == "" {
		b.logger.Debug("failed to obtain password for Basic auth")
		return false, nil
	}

	b.username = username
	b.password = password
	req.SetBasicAuth(b.username, b.password)
	return true, nil
}

// AuthCompleted finishes authentication for the Basic auth responder.
func (b *Basic) AuthCompleted(challenge string, resp *http.Response) (result bool, err error) {
	if challenge == "" {
		return true, nil
	}
	return false, errUnexpectedAuthenticateHeader
}
