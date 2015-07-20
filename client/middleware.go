package client

import (
	"net/http"

	"github.com/docker/engine-api/client/authn"
)

func (cli *Client) doWithMiddlewares(d func(*http.Request) (*http.Response, error)) func(*http.Request) (*http.Response, error) {
	middlewares := []func(func(*http.Request) (*http.Response, error)) func(*http.Request) (*http.Response, error){
		cli.cookieMiddleware,
		authn.Middleware(cli.logger, cli.authers...),
	}
	for _, m := range middlewares {
		d = m(d)
	}
	return d
}
