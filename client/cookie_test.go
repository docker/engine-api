package client

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"

	"github.com/docker/engine-api/client/transport"
)

func cookieMock(req *http.Request) (*http.Response, error) {
	c, err := req.Cookie("engine-api-test-cookie")
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("unable to find engine-api-test-cookie in the request")
	}
	return infoMock(req)
}

func TestClientWithCookieMiddleware(t *testing.T) {
	cookie := &http.Cookie{
		Name:  "engine-api-test-cookie",
		Value: "true",
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	u, err := url.Parse("http://localhost/info")
	if err != nil {
		t.Fatal(err)
	}
	jar.SetCookies(u, []*http.Cookie{cookie})

	client := &Client{
		addr:      "localhost",
		transport: transport.NewMockClient(nil, cookieMock),
	}
	client.AddMiddlewares(transport.NewCookieJarMiddleware(jar))

	_, err = client.Info()
	if err != nil {
		t.Fatal(err)
	}
}
