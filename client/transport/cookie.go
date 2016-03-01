package transport

import (
	"fmt"
	"net/http"
)

type cookieJar struct {
	jar  http.CookieJar
	next Sender
}

// NewCookieJarMiddleware creates a new middleware that injects cookies to requests.
func NewCookieJarMiddleware(jar http.CookieJar) func(next Sender) Sender {
	return func(next Sender) Sender {
		return &cookieJar{
			jar:  jar,
			next: next,
		}
	}
}

// Do sends a request to the next sender and injects cookies for the given URL.
func (c *cookieJar) Do(req *http.Request) (*http.Response, error) {
	fmt.Println(req.URL)
	for _, cookie := range c.jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}
	resp, err := c.next.Do(req)
	if resp != nil {
		if cookies := resp.Cookies(); cookies != nil {
			c.jar.SetCookies(req.URL, cookies)
		}
	}
	return resp, err
}
