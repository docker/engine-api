package authn

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// NegotiateAuther is an interface which a caller may provide for telling us if
// we should attempt Negotiate authentication with a server.
type NegotiateAuther interface {
	GetNegotiateAuth() bool
}

// Middleware returns a function which wraps the passed-in Do()-style function,
// handling any "unauthorized" errors which it returns by retrying the same
// request with authentication.
func Middleware(logger Logger, authers ...interface{}) func(func(req *http.Request) (resp *http.Response, err error)) func(req *http.Request) (resp *http.Response, err error) {
	authResponders := createAuthResponders(logger, authers)
	return func(doer func(req *http.Request) (resp *http.Response, err error)) func(req *http.Request) (resp *http.Response, err error) {
		return func(req *http.Request) (resp *http.Response, err error) {
			// We may have to issue the request multiple times, so
			// we need to be able to rewind and recover everything
			// we've sent.
			rewinder := makeRewinder(req.Body)
			req.Body = rewinder
			resp, err = doer(req)
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
					responder, ok := authResponders[strings.ToLower(tokens[0])]
					if !ok {
						logger.Debug(fmt.Sprintf("no support for authentication scheme \"%s\"", tokens[0]))
						continue
					}
					retryWithUpdatedAuthn, err = responder.authRespond(challenge, req)
					if retryWithUpdatedAuthn {
						logger.Debug(fmt.Sprintf("handler for \"%s\" produced data", tokens[0]))
						scheme = strings.ToLower(tokens[0])
						break
					}
					if err != nil {
						logger.Debug(fmt.Sprintf("%v. handler for \"%s\" failed to produce data", err, tokens[0]))
					} else {
						logger.Debug(fmt.Sprintf("handler for \"%s\" failed to produce data", tokens[0]))
					}
				}
				if len(ah) == 0 {
					if triedAuthnPreviously {
						err = fmt.Errorf("Failed to authenticate to docker daemon")
					} else {
						err = errors.New("Failed to authenticate to docker daemon; server offered no authentication methods")
					}
					break
				} else if err != nil {
					err = fmt.Errorf("%v. Failed to authenticate to docker daemon", err)
					break
				} else if !retryWithUpdatedAuthn {
					err = errors.New("Unable to attempt to authenticate to docker daemon")
					break
				} else {
					ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					rewinder.Rewind()
					req.Body = ioutil.NopCloser(rewinder)
					resp, err = doer(req)
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
					responder := authResponders[scheme]
					completed, err = responder.authCompleted(strings.Join(tokens, " "), resp)
					if completed {
						logger.Debug(fmt.Sprintf("handler for \"%s\" succeeded", scheme))
					} else {
						logger.Debug(fmt.Sprintf("handler for \"%s\" failed", scheme))
					}
				} else if len(ah) == 0 {
					logger.Debug("No authentication header in final server response")
				} else if err != nil {
					err = fmt.Errorf("%v. Unable to authenticate docker daemon", err)
				} else if !completed {
					err = fmt.Errorf("Unable to authenticate docker daemon")
				}
			}
			return resp, err
		}
	}
}
