package authn

import (
	"errors"
	"net/http"

	"github.com/docker/engine-api/client/logger"
)

// AuthResponder is an interface that wraps the scheme,
// authRespond, and authCompleted methods.
//
// At initialization time, an implementation of authResponder should register
// itself by calling registerAuthResponder.
type AuthResponder interface {
	// AuthCompleted, given a (possibly empty) WWW-Authenticate header and
	// a successful response, should decide if the server's reply should be
	// accepted.
	AuthCompleted(challenge string, resp *http.Response) (bool, error)
	// AuthRespond, given the authentication header value associated with
	// the scheme that it implements, can decide if the request should be
	// retried.  If it returns true, then the request is retransmitted to
	// the server, presumably because it has added an authentication header
	// which it believes the server will accept.
	AuthRespond(challenge string, req *http.Request) (bool, error)
	// Scheme should return the name of the authorization scheme for which
	// the responder should be called.
	Scheme() string
	// SetLogger allows a caller to set a logger for the AuthResponder.
	SetLogger(logger.Logger)
}

var errUnexpectedAuthenticateHeader = errors.New("Error: unexpected WWW-Authenticate header in server response")
