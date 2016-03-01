package authn

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/docker/engine-api/client/logger"
	"github.com/docker/engine-api/client/middleware"
	"github.com/docker/engine-api/client/transport"
)

type authMiddleware struct {
	authResponders map[string]AuthResponder
	logger         logger.Logger
	next           transport.Sender
}

// NewAuthResponderMiddleware returns a function which wraps the passed-in Do()-style function,
// handling any "unauthorized" errors which it returns by retrying the same
// request with authentication.
func NewAuthResponderMiddleware(logger logger.Logger, auth ...AuthResponder) middleware.Middleware {
	responders := make(map[string]AuthResponder)
	for _, a := range auth {
		a.SetLogger(logger)
		responders[strings.ToLower(a.Scheme())] = a
	}

	return func(next transport.Sender) transport.Sender {
		return authMiddleware{
			authResponders: responders,
			logger:         logger,
			next:           next,
		}
	}
}

func (a authMiddleware) Do(req *http.Request) (resp *http.Response, err error) {
	// We may have to issue the request multiple times, so
	// we need to be able to rewind and recover everything
	// we've sent.
	var body bytes.Buffer
	if req.Body != nil {
		io.Copy(&body, req.Body)
		if closer, ok := req.Body.(io.Closer); ok {
			closer.Close()
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))
	}
	resp, err = a.next.Do(req)
	// If we previously tried to authenticate, or this
	// isn't an authentication-required error, we're done.
	if req.Header.Get("Authorization") != "" || err != nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}
	// Handle Unauthorized errors by attempting to
	// authenticate, possibly doing so over multiple round
	// trips.
	scheme := ""
	reqheader := http.CanonicalHeaderKey("Authorization")
	respheader := http.CanonicalHeaderKey("WWW-Authenticate")

	for err == nil && resp.StatusCode == http.StatusUnauthorized {
		authnHeaders := req.Header[reqheader]
		triedAuthnPreviously := authnHeaders != nil && len(authnHeaders) > 0
		retryWithUpdatedAuthn := false
		ah := resp.Header[respheader]
		for _, challenge := range ah {
			tokens := strings.Split(strings.Replace(challenge, "\t", " ", -1), " ")
			responder, ok := a.authResponders[strings.ToLower(tokens[0])]
			if !ok {
				a.logger.Debugf("no support for authentication scheme \"%s\"", tokens[0])
				continue
			}
			retryWithUpdatedAuthn, err = responder.AuthRespond(challenge, req)
			if retryWithUpdatedAuthn {
				a.logger.Debugf("handler for \"%s\" produced data", tokens[0])
				scheme = strings.ToLower(tokens[0])
				break
			}
			if err != nil {
				a.logger.Debugf("%v. handler for \"%s\" failed to produce data", err, tokens[0])
			} else {
				a.logger.Debugf("handler for \"%s\" failed to produce data", tokens[0])
			}
		}

		if len(ah) == 0 {
			if triedAuthnPreviously {
				err = fmt.Errorf("Failed to authenticate to docker daemon")
			} else {
				err = fmt.Errorf("Failed to authenticate to docker daemon; server offered no authentication methods")
			}
			break
		} else if err != nil {
			err = fmt.Errorf("%v. Failed to authenticate to docker daemon", err)
			break
		} else if !retryWithUpdatedAuthn {
			err = fmt.Errorf("Unable to attempt to authenticate to docker daemon")
			break
		} else {
			ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if req.Body != nil {
				req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))
			}
			resp, err = a.next.Do(req)
		}
	}

	if err == nil && resp.StatusCode != http.StatusUnauthorized {
		completed := false
		tokens := []string{}
		ah := resp.Header[respheader]
		for _, challenge := range ah {
			tokens = strings.Split(strings.Replace(challenge, "\t", " ", -1), " ")
			if strings.ToLower(tokens[0]) == scheme {
				break
			}
		}

		if len(tokens) == 0 || strings.ToLower(tokens[0]) == scheme {
			responder := a.authResponders[scheme]
			completed, err = responder.AuthCompleted(strings.Join(tokens, " "), resp)
			if completed {
				a.logger.Debugf("handler for \"%s\" succeeded", scheme)
			} else {
				a.logger.Debugf("handler for \"%s\" failed", scheme)
			}
		} else if len(ah) == 0 {
			a.logger.Debug("No authentication header in final server response")
		} else if err != nil {
			err = fmt.Errorf("%v. Unable to authenticate docker daemon", err)
		} else if !completed {
			err = fmt.Errorf("Unable to authenticate docker daemon")
		}
	}

	return resp, err
}
