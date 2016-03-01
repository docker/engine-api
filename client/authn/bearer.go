package authn

import (
	"fmt"
	"net/http"

	"github.com/docker/engine-api/client/logger"
)

// BearerAuthCallback is an interface which a caller may provide for obtaining a
// token to use in attempting bearer authentication with a server.
type BearerAuthCallback func(challenge string) (token string, err error)

// bearer is an AuthResponder that handles bearer authentication.
type bearer struct {
	logger   logger.Logger
	token    string
	callback BearerAuthCallback
}

// NewBearerAuth creates a Bearer auth responder with a callback
func NewBearerAuth(callback BearerAuthCallback) AuthResponder {
	return &bearer{
		logger:   logger.Silent{},
		callback: callback,
	}
}

// SetLogger sets the logger for the bearer auth responder.
func (b *bearer) SetLogger(l logger.Logger) {
	b.logger = l
}

// Scheme returns the scheme the bearer auth responder handles.
func (b *bearer) Scheme() string {
	return "Bearer"
}

// AuthRespond handles authentication for the bearer auth responder.
func (b *bearer) AuthRespond(challenge string, req *http.Request) (result bool, err error) {
	if b.token != "" {
		b.logger.Debug("using previously-supplied Bearer token")
		req.Header.Add("Authorization", b.header())
		return true, nil
	}
	if b.callback == nil {
		b.logger.Debug("failed to obtain token for Bearer auth")
		return false, nil
	}
	token, err := b.callback(challenge)
	if err != nil {
		return false, err
	}
	if token == "" {
		b.logger.Debug("Bearer token not supplied")
		return false, nil
	}
	b.token = token
	req.Header.Add("Authorization", b.header())
	return true, nil
}

// AuthCompleted finishes authentication for the bearer auth responder.
func (b *bearer) AuthCompleted(challenge string, resp *http.Response) (result bool, err error) {
	if challenge == "" {
		return true, nil
	}
	return false, errUnexpectedAuthenticateHeader
}

func (b *bearer) header() string {
	return fmt.Sprintf("%s %s", b.Scheme(), b.token)
}
